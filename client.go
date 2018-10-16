package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	currentChannel *Channel
	currentUser    *User
	welcomeMessage = `
	W E L C O M E   T O
_____  ____   _____ _    _       _______
/ ____|/ __ \ / ____| |  | |   /\|__   __|
| |  __| |  | | |    | |__| |  /  \  | |
| | |_ | |  | | |    |  __  | / /\ \ | |
| |__| | |__| | |____| |  | |/ ____ \| |
\_____|\____/ \_____|_|  |_/_/    \_|_|

Type any message you like and hit <return> to send to the channel.
Type /help<return> for a list of commands.
Type /exit<return> to leave chat.

`
)

// Client contains the client connection and the channel for dispatching
// messages back to the server for storage and dispatch to other connected clients
type Client struct {
	Channel    *Channel
	Conn       net.Conn
	Context    *Context
	Dispatcher chan *Message
	Mux        *sync.Mutex
	Receiver   chan *Message
	Scanner    *bufio.Scanner
	User       *User
	Writer     *bufio.Writer
}

func newClient(conn net.Conn, context *Context, dispatcher chan *Message) *Client {
	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanLines)
	return &Client{
		Channel:    newChannel("general"),
		Conn:       conn,
		Context:    context,
		Dispatcher: dispatcher,
		Mux:        &sync.Mutex{},
		Receiver:   make(chan *Message),
		Scanner:    scanner,
		User:       &User{"anonymous"},
		Writer:     bufio.NewWriter(conn),
	}
}

func (client *Client) connected() {
	client.Context.Logger.Debugf("Client connected from %s", client.Conn.RemoteAddr())
	client.writeString(welcomeMessage)
	client.login()
	go client.listen()
	go client.receive()
}

func (client *Client) dispatchMessage(text string) {
	client.Dispatcher <- &Message{
		Channel: client.Channel,
		Message: text,
		Time:    time.Now().UTC(),
		User:    client.User,
	}
}

func getWord(line string, index int) string {
	words := strings.Fields(line)
	return strings.TrimSpace(words[index])
}

func (client *Client) joinChannel(channel *Channel) {
	client.Mux.Lock()
	defer client.Mux.Unlock()
	client.Channel = channel
}

func (client *Client) listen() {
	defer func() {
		client.Context.Logger.Debugf("Client disconnected from %s", client.Conn.RemoteAddr())
		client.Conn.Close()
	}()
	for {
		client.writePrompt()

		if ok := client.Scanner.Scan(); !ok {
			break
		}
		channelMessage := client.Scanner.Text()
		client.dispatchMessage(channelMessage)
	}
}

func (client *Client) login() {
	defer func() {
		if r := recover(); r != nil {
			client.Context.Logger.Errorf("Unable to login [%T]: %v", r, r)
		}
	}()

	if username, err := client.prompt("Username: "); nil != err {
		panic(err)
	} else {
		client.User.Username = getWord(username, 0)
	}

	client.dispatchMessage(fmt.Sprintf("/join %s", client.Channel.Name))
}

func (client *Client) prompt(question string) (line string, err error) {
	if err = client.writeString(question); nil != err {
		panic(err)
	}

	ok := client.Scanner.Scan()
	if !ok {
		panic(client.Scanner.Err())
	}
	line = client.Scanner.Text()
	return
}

func (client *Client) receive() {
	for {
		select {
		case message := <-client.Receiver:
			client.writeString("\r")
			client.writeLine(message.String())
		}
		client.writePrompt()
	}
}

func (client *Client) writeLine(line string) (err error) {
	return client.writeString(fmt.Sprintf("%s\n", line))
}

func (client *Client) writePrompt() {
	client.writeString("$ ")
}

func (client *Client) writeString(line string) (err error) {
	if _, err = client.Writer.WriteString(line); nil != err {
		return err
	}
	if err = client.Writer.Flush(); nil != err {
		return err
	}
	return nil
}
