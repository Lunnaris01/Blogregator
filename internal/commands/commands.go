package commands

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"html"
	"io"
	"net/http"
	"os"
	"time"
	"database/sql"
	"github.com/google/uuid"
	"strconv"
)

// Define state for handlers
type State struct {
	Db     *database.Queries
	Config *config.Config
}

// Define the command struct
type Command struct {
	Name string
	Args []string
}

type Commands struct {
	handlers map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.handlers == nil {
		c.handlers = make(map[string]func(*State, Command) error)
	}
	c.handlers[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, exists := c.handlers[cmd.Name]; exists {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}

// Login handler example
func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("a username is required")
	}
	username := cmd.Args[0]
	new_context := context.Background()
	_, err := s.Db.GetUser(new_context, username)
	if err != nil {
		// User exists!
		fmt.Printf("User %s does not exist\n", cmd.Args[0])
		os.Exit(1)
	}

	err = config.SetUser(s.Config, username) // Set the current user
	if err != nil {
		return err
	}

	fmt.Printf("User %s is now logged in\n", username)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("a name is required")
	}

	new_context := context.Background()
	_, err := s.Db.GetUser(new_context, cmd.Args[0])
	if err == nil {
		// User exists!
		fmt.Printf("User %s already exists\n", cmd.Args[0])
		os.Exit(1)

	}
	new_user_params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	new_user, err := s.Db.CreateUser(new_context, new_user_params)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	err = config.SetUser(s.Config, new_user.Name)
	if err != nil {
		return fmt.Errorf("failed to set current user: %v", err)
	}
	fmt.Printf("User %s added succesfully\n", new_user.Name)
	return nil

}

func HandlerPrintUsers(s *State, cmd Command) error {
	new_context := context.Background()
	users, err := s.Db.GetUsers(new_context)
	if err != nil {
		fmt.Printf("Error when trying to get Users")
		os.Exit(1)
	}
	for _, u_name := range users {
		if s.Config.CurrentUserName == u_name {
			fmt.Printf("* %s (current)\n", u_name)
		} else {
			fmt.Printf("* %s\n", u_name)
		}
	}
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	new_context := context.Background()
	err := s.Db.ResetUsers(new_context)
	err = s.Db.ResetFeeds(new_context)
	if err != nil {
		fmt.Printf("Error when resetting database")
		os.Exit(1)
	}
	fmt.Println("Resetting Database - Removing all Entries!")
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("A valid time duration is needed!")
	}
	time_as_int, err  := strconv.Atoi(cmd.Args[0])
	if err != nil || time_as_int<15 || time_as_int>60 {
		return errors.New("Choose an update duration between 15 and 60 seconds!")	
	}
	ticker := time.NewTicker(time.Second)
	counter := 0
	fmt.Printf("Collecting feeds every %d seconds\n", time_as_int)
	for ; ; <- ticker.C{
		if counter < time_as_int{
			counter++
			if counter %5 == 0{
				fmt.Printf("Next update in %d Seconds\n",time_as_int-counter)
			}
			continue
		}
		err := scrapeFeeds(s, cmd)
		if err != nil {
			return fmt.Errorf("Failure during update: %v", err)
		}
		counter= 0
	}
	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("Feed needs a Name!")
	} else if len(cmd.Args) < 2 {
		return errors.New("Need URL to create feed!")
	}

	new_context := context.Background()
	new_feed_params := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	new_feed, err := s.Db.AddFeed(new_context, new_feed_params)
	if err != nil {
		return fmt.Errorf("failed to create feed: %v", err)
	}
	fmt.Printf("Feed %s added succesfully\n", new_feed.Name)
	err = followFeed(user.ID, new_feed.ID, s)
	if err != nil {
		return err
	}

	return nil

}

func HandlerPrintFeeds(s *State, cmd Command) error {
	new_context := context.Background()
	feeds, err := s.Db.GetFeeds(new_context)
	if err != nil {
		fmt.Printf("Error when trying to grab the Feeds")
		os.Exit(1)
	}
	for _, feed := range feeds {
		fmt.Printf("  -------------------  \n")
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("Added by User: %s\n", feed.Username)
		fmt.Printf("  -------------------  \n")

	}
	return nil

}

func HandlerFollowFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("Use the URL of the feed to follow!")
	}

	new_context := context.Background()
	feed, err := s.Db.GetFeedByURL(new_context, cmd.Args[0])
	if err != nil {
		fmt.Printf("Error when trying to grab the Feed")
		os.Exit(1)
	}

	err = followFeed(user.ID, feed.ID, s)
	if err != nil {
		return err
	}
	return nil
}

func HandlerUnFollowFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("Use the URL of the feed to unfollow!")
	}

	unfollow_params := database.UnfollowFeedParams{
		Url:	cmd.Args[0],
		UserID:	user.ID,
	}

	new_context := context.Background()
	feedname, err := s.Db.UnfollowFeed(new_context, unfollow_params)
	if err != nil {
		fmt.Printf("Error when trying to unfollow: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s is no longer following %s\n", user.Name,feedname)
	return nil
}


func HandlerFollowing(s *State, cmd Command, user database.User) error {

	followed_feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Printf("Failed to fetch Feeds")
		os.Exit(1)
	}
	fmt.Printf("User %s is following these feeds:\n", user.Name)
	for _, feed := range followed_feeds {
		fmt.Printf("    *%s\n", feed.Feedname)
	}
	return nil
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	rssfeed := RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &rssfeed, err
	}
	req.Header.Set("User-Agent", "gator")
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	res, err := netClient.Do(req)
	if err != nil {
		return &rssfeed, err
	}
	defer res.Body.Close()
	data_body, err := io.ReadAll(res.Body)
	if err != nil {

		return &RSSFeed{}, fmt.Errorf("Failed to read body!")
	}

	err = xml.Unmarshal(data_body, &rssfeed)
	if err != nil {

		return &RSSFeed{}, fmt.Errorf("Failed to Unmarshal body!")
	}
	rssfeed.Channel.Title = html.UnescapeString(rssfeed.Channel.Title)
	rssfeed.Channel.Description = html.UnescapeString(rssfeed.Channel.Description)
	for i, item := range rssfeed.Channel.Item {
		rssfeed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		rssfeed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
	return &rssfeed, nil
}

func followFeed(user_id uuid.UUID, feed_id uuid.UUID, s *State) error {

	followFeedParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user_id,
		FeedID:    feed_id,
	}
	newFeedFollow, err := s.Db.CreateFeedFollow(context.Background(), followFeedParams)
	if err != nil {
		return fmt.Errorf("Failed to follow Feed with Error: %v\n", err)
	}
	fmt.Printf("%s is now following %s\n", newFeedFollow.UserName, newFeedFollow.FeedName)
	return nil
}

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error{

	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	

	}

}

func scrapeFeeds(s *State, cmd Command) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	fetched_feed,err := fetchFeed(context.Background(),feed.Url)
	if err != nil {
		return err
	}
	err = s.Db.MarkFeedFetched(
		context.Background(),
		database.MarkFeedFetchedParams{
			ID:				feed.ID,
			UpdatedAt:	time.Now(),
			LastFetchedAt: sql.NullTime{
				Time: 	time.Now(),
				Valid: 	true,
				},
			})
	if err != nil {
		return err
	}
	
	printFeed(fetched_feed)
	fmt.Println("Feed Updated")
	return nil
}

func printFeed(feed *RSSFeed){
	fmt.Printf("Title: %s\n",feed.Channel.Title)
	for _, item := range feed.Channel.Item{
		fmt.Printf("Title: %s\n",item.Title)
		fmt.Printf("Description: %s\n",item.Description)
	}
}