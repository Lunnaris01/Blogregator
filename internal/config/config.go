package config

import (
	"os"
	"encoding/json"
)


type Config struct{
	Db_url string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error){
	configDir, err := os.UserHomeDir()
	if err != nil{
		return Config{},err
	}
	configDir = configDir + "/.gatorconfig.json"
	dat, err := os.ReadFile(configDir)
	if err != nil{
		return Config{},err
	}
	newConfig := Config{}
	err = json.Unmarshal(dat,&newConfig)
	if err != nil{
		return Config{},err
	}
	return newConfig, nil
	
}

func SetUser(c Config, userName string) error {
	c.CurrentUserName = userName
	fileContent, err := json.Marshal(c)
	if err != nil{
		return err
	}
	configDir, err := os.UserHomeDir()
	if err != nil{
		return err
	}
	configDir = configDir + "/.gatorconfig.json"
	err = os.WriteFile(configDir,fileContent,0644)
	return nil
}