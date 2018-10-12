package main

import (
	"fmt"

	telnet "github.com/reiver/go-telnet"
)

func startServer(config *chatServerConfig) {
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

	if err := server.ListenAndServe(); err != nil {
		printErrorAndExit(err, -1)
	}
	fmt.Printf("Server listening on %v", config.Port)
}
