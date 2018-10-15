package main

import (
	"bufio"
	"fmt"
	"net"
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
	Channel *Channel
	Conn    net.Conn
	Context *Context
	Scanner *bufio.Scanner
	User    *User
	Writer  *bufio.Writer
}

func (client *Client) connected() {
	client.Context.Logger.Debugf("Client connected from %s", client.Conn.RemoteAddr())
	client.Writer = bufio.NewWriter(client.Conn)
	client.Scanner = bufio.NewScanner(client.Conn)
	client.Scanner.Split(bufio.ScanLines)

	client.writeString(welcomeMessage)
}

func (client *Client) listen() {
func (client *Client) writeLine(line string) (err error) {
	return client.writeString(fmt.Sprintf("%s\n", line))
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
