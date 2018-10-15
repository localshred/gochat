package main

import (
	"fmt"
	"net"
)

// Context is a struct that can be passed to clients to get access to main server resources
type Context struct {
	Config *chatServerConfig
	Logger *telnetLogger
}

// TelnetServer represents a wrapped telnet server
type TelnetServer struct {
	Clients  []*Client
	Channels []*Channel
	Context  *Context
}

func (telnetServer *TelnetServer) start() {
	context := telnetServer.Context
	config := context.Config
	address := fmt.Sprintf("%s:%v", config.Host, config.Port)
	logger, logFile, err := createLogger(config.LogFile)
	if err != nil {
		printErrorAndExit(err, -1)
	}
	defer logFile.Close()
	context.Logger = logger

	listener, err := net.Listen("tcp", address)
	if nil != err {
		printErrorAndExit(err, -1)
	}
	defer listener.Close()

	listeningMessage := fmt.Sprintf("Server listening on %s\n", address)
	fmt.Print(listeningMessage)
	context.Logger.Debug(listeningMessage)
	for {
		telnetServer.acceptConnection(listener)
	}
}

func (telnetServer *TelnetServer) acceptConnection(listener net.Listener) {
	conn, err := listener.Accept()
	if nil != err {
		printErrorAndExit(err, -1)
	}

	client := &Client{
		Channel: &Channel{"general"},
		Conn:    conn,
		Context: telnetServer.Context,
		User:    &User{"anonymous"},
	}
	telnetServer.Clients = append(telnetServer.Clients, client)
	client.connected()
	go client.listen()
}
