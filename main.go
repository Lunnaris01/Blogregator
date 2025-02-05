package main

import (
	"fmt"
	"gator/internal/config"
	"log"
)

func main(){
	fmt.Println("Gator is starting up!")
	cfg ,err := config.Read()
	if err != nil{
		fmt.Println(err)
		log.Fatal("Couldnt read config")
	}
	fmt.Println(cfg)
	err = config.SetUser(cfg,"Julian")
	if err != nil{
		fmt.Println(err)
		log.Fatal("Couldnt write config")
	}
	cfg ,err = config.Read()
	if err != nil{
		fmt.Println(err)
		log.Fatal("Couldnt read config")
	}
	fmt.Println(cfg)


}
