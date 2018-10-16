package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

// Context is a struct that can be passed to clients to get access to main server resources
type Context struct {
	Config *chatServerConfig
	Logger *telnetLogger
}

// Server represents a wrapped telnet server
type Server struct {
	Clients    []*Client
	Channels   map[string]*Channel
	Context    *Context
	Dispatcher chan *Message
	Mux        *sync.Mutex
}

func newServer(config *chatServerConfig) *Server {
	return &Server{
		Context: &Context{Config: config},
		Mux:     &sync.Mutex{},
	}
}

func (server *Server) acceptConnection(listener net.Listener) {
	conn, err := listener.Accept()
	if nil != err {
		printErrorAndExit(err, -1)
	}
	server.Mux.Lock()
	defer server.Mux.Unlock()

	client := newClient(conn, server.Context, server.Dispatcher)
	server.Clients = append(server.Clients, client)
	server.Context.Logger.Debugf("Connected Clients: %v", len(server.Clients))
	go client.connected()
}

func (server *Server) commandJoin(channelName string, message *Message) {
	channelToJoin := server.findOrCreateChannel(channelName)

	if channelName != message.Channel.Name {
		channelToLeave := message.Channel
		leaveMessage := channelToLeave.userLeft(message.User)
		server.sendToClients(leaveMessage)

		for _, client := range server.Clients {
			if client.Channel == channelToLeave {
				client.joinChannel(channelToJoin)
			}
		}
	}

	joinMessage := channelToJoin.userJoined(message.User)
	server.sendToClients(joinMessage)
}

func (server *Server) findOrCreateChannel(channelName string) *Channel {
	// TODO mutex?
	channel, ok := server.Channels[channelName]
	if !ok {
		channel = newChannel(channelName)
		server.Channels[channelName] = channel
	}
	return channel
}

func (server *Server) receiveFromClients() {
	for {
		select {
		case message := <-server.Dispatcher:
			fields := strings.Fields(message.Message)
			switch fields[0] {
			case "/join":
				server.commandJoin(fields[1], message)
			default:
				server.sendToClients(message)
			}
		}
	}
}

func (server *Server) sendToClients(message *Message) {
	message.Channel.appendMessage(server.Context, message)
	server.Mux.Lock()
	defer server.Mux.Unlock()
	for _, client := range server.Clients {
		if client.Channel.Name == message.Channel.Name {
			client.Receiver <- message
		}
	}
}

func (server *Server) start() {
	context := server.Context
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

	server.Channels = map[string]*Channel{}
	server.findOrCreateChannel("general")

	server.Dispatcher = make(chan *Message)
	go server.receiveFromClients()
	for {
		server.acceptConnection(listener)
	}
}
