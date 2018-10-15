package main

import (
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
Type help<return> for a list of commands.
Type exit<return> to leave chat.

`
)

// Client contains the client connection and the channel for dispatching
// messages back to the server for storage and dispatch to other connected clients
type Client struct {
	Channel *Channel
	Conn    net.Conn
	User    *User
}

func (client *Client) connected() {
	client.Conn.Write([]byte(welcomeMessage))
}

func (client *Client) listen() {

}
