package main

import (
	"fmt"
	"gator/internal/config"
	"gator/internal/commands"
	"gator/internal/database"
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
	cmds.Register("addfeed",commands.HandlerAddFeed)
	cmds.Register("feeds" , commands.HandlerPrintFeeds)

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
