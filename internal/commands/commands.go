package commands

import (
	"errors"
	"gator/internal/config"
	"gator/internal/database"
	"fmt"
	"context"
	"time"
	"github.com/google/uuid"
	"os"
	"net/http"
	"encoding/xml"
	"io"
	"html"
)


// Define state for handlers
type State struct {
	Db *database.Queries
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
		c.handlers = make(map[string]func(*State,Command) error)
	}
	c.handlers[name] = f
}

func (c *Commands) Run (s *State, cmd Command) error {
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
	if err != nil{
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
	if err == nil{
		// User exists!
		fmt.Printf("User %s already exists\n", cmd.Args[0])
		os.Exit(1)
		
	}
	new_user_params := database.CreateUserParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		cmd.Args[0],
	}

	new_user, err := s.Db.CreateUser(new_context,new_user_params)
	if err != nil {
        return fmt.Errorf("failed to create user: %v", err)
	}
	err = config.SetUser(s.Config, new_user.Name)
	if err != nil {
        return fmt.Errorf("failed to set current user: %v", err)
	}
	fmt.Printf("User %s added succesfully\n",new_user.Name)
	return nil

}

func HandlerPrintUsers(s *State, cmd Command) error {
	new_context := context.Background()
	users,err := s.Db.GetUsers(new_context)
	if err != nil {
		fmt.Printf("Error when trying to get Users")
		os.Exit(1)
	}
	for _, u_name := range users {
		if s.Config.CurrentUserName == u_name{
			fmt.Printf("* %s (current)\n",u_name)
		} else {
			fmt.Printf("* %s\n",u_name)
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
	new_context := context.Background()
	cheat_url := "https://www.wagslane.dev/index.xml"
	rssfeed, err := fetchFeed(new_context,cheat_url)
	if err != nil{
		return fmt.Errorf("Failed to fetch content: %v", err)
	}
	fmt.Printf("%#v\n",rssfeed)
	return nil
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	rssfeed := RSSFeed{}
	req,err := http.NewRequestWithContext(ctx,"GET",feedURL,nil)
	if err != nil{
		return &rssfeed,err
	}
	req.Header.Set("User-Agent","gator")
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	res, err := netClient.Do(req)
	if err != nil{
		return &rssfeed,err
	}
	defer res.Body.Close()
	data_body, err := io.ReadAll(res.Body)
	if err != nil {
		
		return &RSSFeed{},fmt.Errorf("Failed to read body!")
	}

	err = xml.Unmarshal(data_body,&rssfeed)
	if err != nil {
		
		return &RSSFeed{},fmt.Errorf("Failed to Unmarshal body!")
	}
	rssfeed.Channel.Title = html.UnescapeString(rssfeed.Channel.Title)
	rssfeed.Channel.Description = html.UnescapeString(rssfeed.Channel.Description)
	for i,item := range rssfeed.Channel.Item {
		rssfeed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		rssfeed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
	return &rssfeed,nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("Feed needs a Name!")
	} else if len(cmd.Args) < 2{
		return errors.New("Need URL to create feed!")
	}

	new_context := context.Background()
	current_user := s.Config.CurrentUserName
	user_id, err := s.Db.GetUserIdByName(new_context, current_user)
	if err != nil{
		fmt.Printf("Failed to fetch id for current User: %s\n", current_user)
		os.Exit(1)
	}

	new_feed_params := database.AddFeedParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		cmd.Args[0],
		Url:		cmd.Args[1],
		UserID:		user_id,
	}

	new_feed, err := s.Db.AddFeed(new_context,new_feed_params)
	if err != nil {
        return fmt.Errorf("failed to create feed: %v", err)
	}
	fmt.Printf("Feed %s added succesfully\n",new_feed.Name)
	return nil

}


func HandlerPrintFeeds(s *State, cmd Command) error {
	new_context := context.Background()
	feeds,err := s.Db.GetFeeds(new_context)
	if err != nil {
		fmt.Printf("Error when trying to grab the Feeds")
		os.Exit(1)
	}
	for _, feed := range feeds {
		fmt.Printf("  -------------------  \n")
		fmt.Printf("Name: %s\n",feed.Name)
		fmt.Printf("URL: %s\n",feed.Url)
		fmt.Printf("Added by User: %s\n",feed.Username)
		fmt.Printf("  -------------------  \n")

	}
	return nil

}