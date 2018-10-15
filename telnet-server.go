package main

import (
	"fmt"
	"net"
)

// TelnetServer represents a wrapped telnet server
type TelnetServer struct {
	Clients  []*Client
	Channels []*Channel
	Config   *chatServerConfig
	Logger   *telnetLogger
}

func (telnetServer *TelnetServer) start() {
	config := telnetServer.Config
	address := fmt.Sprintf("%s:%v", config.Host, config.Port)
	logger, logFile, err := createLogger(config.LogFile)
	if err != nil {
		printErrorAndExit(err, -1)
	}
	defer logFile.Close()
	telnetServer.Logger = logger

	listener, err := net.Listen("tcp", address)
	if nil != err {
		printErrorAndExit(err, -1)
	}
	defer listener.Close()

	listeningMessage := fmt.Sprintf("Server listening on %s\n", address)
	fmt.Print(listeningMessage)
	telnetServer.Logger.Debug(listeningMessage)
	for {
		telnetServer.acceptConnection(listener)
	}
}

func (telnetServer *TelnetServer) acceptConnection(listener net.Listener) {
	conn, err := listener.Accept()
	if nil != err {
		printErrorAndExit(err, -1)
	}
	defer func() {
		telnetServer.Logger.Debugf("Client disconnected from %s", conn.RemoteAddr())
		conn.Close()
	}()

	telnetServer.Logger.Debugf("Client connected from %s", conn.RemoteAddr())

	client := &Client{
		Channel: &Channel{"general"},
		Conn:    conn,
		User:    &User{"anonymous"},
	}
	telnetServer.Clients = append(telnetServer.Clients, client)

	client.connected()
	go client.listen()
	// TODO spawn go-routine to read client messages
}
