package main

import (
	"fmt"
	"net"
	"strings"
)

// Context is a struct that can be passed to clients to get access to main server resources
type Context struct {
	Config *chatServerConfig
	Logger *telnetLogger
}

// TelnetServer represents a wrapped telnet server
type TelnetServer struct {
	Clients    []*Client
	Channels   map[string]*Channel
	Context    *Context
	Dispatcher chan *Message
}

func (telnetServer *TelnetServer) acceptConnection(listener net.Listener) {
	conn, err := listener.Accept()
	if nil != err {
		printErrorAndExit(err, -1)
	}

	client := &Client{
		Channel:    newChannel("general"),
		Conn:       conn,
		Context:    telnetServer.Context,
		Dispatcher: telnetServer.Dispatcher,
		Receiver:   make(chan *Message),
		User:       &User{"anonymous"},
	}
	telnetServer.Clients = append(telnetServer.Clients, client)
	telnetServer.Context.Logger.Debugf("Connected Clients: %v", len(telnetServer.Clients))
	go client.connected()
}

func (telnetServer *TelnetServer) commandJoin(channelName string, message *Message) {
	channelToJoin := telnetServer.findOrCreateChannel(channelName)

	if channelName != message.Channel.Name {
		channelToLeave := message.Channel
		leaveMessage := channelToLeave.userLeft(message.User)
		telnetServer.sendToClients(leaveMessage)

		for _, client := range telnetServer.Clients {
			if client.Channel == channelToLeave {
				client.Channel = channelToJoin
			}
		}
	}

	joinMessage := channelToJoin.userJoined(message.User)
	telnetServer.sendToClients(joinMessage)
}

func (telnetServer *TelnetServer) findOrCreateChannel(channelName string) *Channel {
	channel, ok := telnetServer.Channels[channelName]
	if !ok {
		channel = newChannel(channelName)
		telnetServer.Channels[channelName] = channel
	}
	return channel
}

func (telnetServer *TelnetServer) receiveFromClients() {
	for {
		select {
		case message := <-telnetServer.Dispatcher:
			fields := strings.Fields(message.Message)
			switch fields[0] {
			case "/join":
				telnetServer.commandJoin(fields[1], message)
			default:
				telnetServer.sendToClients(message)
			}
		}
	}
}

func (telnetServer *TelnetServer) sendToClients(message *Message) {
	message.Channel.appendMessage(telnetServer.Context, message)
	// TODO lock clients mutex
	for _, client := range telnetServer.Clients {
		if client.Channel.Name == message.Channel.Name {
			client.Receiver <- message
		}
	}
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

	telnetServer.Channels = map[string]*Channel{}
	telnetServer.findOrCreateChannel("general")

	telnetServer.Dispatcher = make(chan *Message)
	go telnetServer.receiveFromClients()
	for {
		telnetServer.acceptConnection(listener)
	}
}
