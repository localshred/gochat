package main

import (
	"fmt"
	"os"

	telnet "github.com/reiver/go-telnet"
)
func main() {
	config := readConfig()
	handler := telnet.EchoHandler
	server := &telnet.Server{
		Addr:    fmt.Sprintf("%s:%v", config.Host, config.Port),
		Handler: handler,
	}

	err := server.ListenAndServe()
	if nil != err {
		fmt.Fprintln(os.Stderr, "Unable to start server: ", err)
		os.Exit(-1)
	}
	fmt.Printf("Server listening on %v", config.Port)
}

type chatServerConfig struct {
	Host    string
	Port    int
	LogFile string
}

func readConfig() *chatServerConfig {
	defaultConfig := &chatServerConfig{
		Host:    "localHost",
		LogFile: "logs/server.log",
		Port:    5555,
	}
	return defaultConfig
}
