package main

import (
	"fmt"
	"time"
)

// Message holds a reference to the channel and the user who posted the message
type Message struct {
	Channel *Channel
	Message string
	Time    time.Time
	User    *User
}

func botMessage(channel *Channel, message string) *Message {
	return &Message{
		Channel: channel,
		Message: message,
		Time:    time.Now().UTC(),
		User:    botUser,
	}
}

func (message *Message) String() string {
	return fmt.Sprintf(
		"[%s] %s: %s",
		message.Time.Local().Format("03:04"),
		message.User.Username,
		message.Message,
	)
}
