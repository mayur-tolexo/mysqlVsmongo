package common

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config model
type Config struct {
	Database struct {
		Mongo struct {
			Hostname   string `json:"hostname"`
			Port       string `json:"port"`
			Database   string `json:"database"`
			Collection string `json:"collection"`
		} `json:"mongo"`
		Mysql struct {
		} `json:"mysql"`
	} `json:"database"`
}

var config *Config

// GetConfig will return the configuration
func GetConfig() *Config {
	if config != nil {
		return config
	}
	file := "config.json"
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
