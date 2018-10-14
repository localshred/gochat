package main

import (
	"fmt"

	telnet "github.com/reiver/go-telnet"
)

// TelnetServer represents a wrapped telnet server
type TelnetServer struct {
	Config *chatServerConfig
}

func (telnetServer *TelnetServer) start() {
	config := telnetServer.Config
	address := fmt.Sprintf("%s:%v", config.Host, config.Port)
	handler := telnet.EchoHandler
	logger, logFile, err := createLogger(config.LogFile)
	if err != nil {
		printErrorAndExit(err, -1)
	}
	defer logFile.Close()

	server := &telnet.Server{
		Addr:    address,
		Handler: handler,
		Logger:  logger,
	}

	fmt.Printf("Server listening on %s\n", address)
	if err := server.ListenAndServe(); err != nil {
		printErrorAndExit(err, -1)
	}
}
