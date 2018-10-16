package main

import (
	"fmt"
	"time"
)

// Channel holds properties related to channels
type Channel struct {
	Name string
}

// Messages are stored as a slice in a map keyed by the channel name
var channelsMessages = map[*Channel][]*Message{}

func newChannel(name string) *Channel {
	return &Channel{
		Name:  name,
	}
}

// TODO lock mutex
func (channel *Channel) appendMessage(context *Context, message string, user *User) (msg *Message) {
	msg = &Message{
		Channel: channel,
		Message: message,
		Time:    time.Now().UTC(),
		User:    user,
	}
	channelsMessages[channel] = append(channelsMessages[channel], msg)
	context.Logger.Debugf("[%s] %s: %s", channel, user.Username, message)
	return
}

func (channel *Channel) String() string {
	return fmt.Sprintf("#%s", channel.Name)
}
