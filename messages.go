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

func (message *Message) String() string {
	return fmt.Sprintf(
		"[%s] %s: %s",
		message.Time.Local().Format("03:04"),
		message.User.Username,
		message.Message,
	)
}
