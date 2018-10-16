package main

import (
	"fmt"
)

// Channel holds properties related to channels
type Channel struct {
	Name     string
	Messages []*Message
	Users    map[string]*User
}

func newChannel(name string) *Channel {
	return &Channel{
		Name:     name,
		Messages: []*Message{},
		Users:    map[string]*User{},
	}
}

func (channel *Channel) appendMessage(context *Context, message *Message) {
	// TODO lock mutex
	channel.Messages = append(channel.Messages, message)
	context.Logger.Debugf("[%s] %s", channel, message)
}

func (channel *Channel) userJoined(user *User) *Message {
	channel.Users[user.Username] = user
	return botMessage(channel, fmt.Sprintf("%s joined %s", user, channel))
}

func (channel *Channel) userLeft(user *User) *Message {
	delete(channel.Users, user.Username)
	return botMessage(channel, fmt.Sprintf("%s left %s", user, channel))
}

func (channel *Channel) String() string {
	return fmt.Sprintf("#%s", channel.Name)
}
