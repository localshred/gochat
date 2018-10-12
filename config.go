package main

import (
	"encoding/json"
	"io/ioutil"
)

var configFile = "./configs/config.json"

type chatServerConfig struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	LogFile string `json:"logFile"`
}

func readConfig() *chatServerConfig {
	config := &chatServerConfig{
		Host:    "localHost",
		LogFile: "logs/server.log",
		Port:    5555,
	}
	rawJSON, err := ioutil.ReadFile(configFile)
	if err != nil {
		printErrorAndExit(err, -1)
	}
	json.Unmarshal([]byte(rawJSON), &config)
	return config
}
