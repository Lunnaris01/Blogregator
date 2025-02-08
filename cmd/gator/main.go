package main

import (
	"fmt"
	"github.com/Lunnaris01/Blogregator/internal/config"
	"github.com/Lunnaris01/Blogregator/internal/commands"
	"github.com/Lunnaris01/Blogregator/internal/database"
	"log"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
)


func main(){
	fmt.Println("Gator is starting up!")
	cfg ,err := config.Read()
	if err != nil{
		log.Fatalf("Error: %v",err)
	}
	appState := commands.State{
		Config: &cfg,
	}

	db, err := sql.Open("postgres",appState.Config.Db_url)
	appState.Db = database.New(db)

	cmds := &commands.Commands{}
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users",commands.HandlerPrintUsers)
	cmds.Register("agg",commands.HandlerAgg)
	cmds.Register("addfeed",commands.MiddlewareLoggedIn(commands.HandlerAddFeed))
	cmds.Register("feeds" , commands.HandlerPrintFeeds)
	cmds.Register("follow" , commands.MiddlewareLoggedIn(commands.HandlerFollowFeed))
	cmds.Register("following" , commands.MiddlewareLoggedIn(commands.HandlerFollowing))
	cmds.Register("unfollow" , commands.MiddlewareLoggedIn(commands.HandlerUnFollowFeed))
	//cmds.Register("scrapefeed", commands.HandlerscrapeFeeds)
	cmds.Register("explore", commands.MiddlewareLoggedIn(commands.HandlerExplorePosts))
	cmds.Register("posts", commands.HandlerPrintPosts)
	args := os.Args[1:]
	if len(args) <1 {
		log.Fatalf("Error: Missing command")
	}

	command := commands.Command{
		Name: args[0],
		Args: args[1:],
	}

	if err := cmds.Run(&appState, command); err != nil {
		log.Fatalf("Error: %v", err)
	}


}
