package main

import (
	"encoding/json"
	"io/ioutil"
)

var configFile = "./configs/config.json"

type chatConfig struct {
	Server     *chatServerConfig `json:"server"`
	HTTPServer *httpServerConfig `json:"http"`
}

type chatServerConfig struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	LogFile string `json:"logFile"`
}

type httpServerConfig struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	LogFile string `json:"logFile"`
}

func readConfig() *chatConfig {
	config := &chatConfig{
		HTTPServer: &httpServerConfig{
			Host:    "localHost",
			LogFile: "logs/server.http.log",
			Port:    5566,
		},
		Server: &chatServerConfig{
			Host:    "localHost",
			LogFile: "logs/server.log",
			Port:    5555,
		},
	}
	rawJSON, err := ioutil.ReadFile(configFile)
	if err != nil {
		printErrorAndExit(err, -1)
	}
	json.Unmarshal([]byte(rawJSON), config)
	return config
}
