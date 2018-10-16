package main

import (
	"fmt"
	"sort"
	"sync"
)

// Channel holds properties related to channels
type Channel struct {
	Name     string           `json:"name"`
	Messages []*Message       `json:"messages"`
	Mux      *sync.Mutex      `json:"-"`
	Users    map[string]*User `json:"users"`
}

func listChannels(channels map[string]*Channel) []string {
	list := []string{}
	for channelName := range channels {
		list = append(list, channelName)
	}
	sort.Strings(list)
	return list
}

func newChannel(name string) *Channel {
	return &Channel{
		Name:     name,
		Messages: []*Message{},
		Mux:      &sync.Mutex{},
		Users:    map[string]*User{},
	}
}

func (channel *Channel) userLeft(user *User) *Message {
	channel.Mux.Lock()
	defer channel.Mux.Unlock()
	delete(channel.Users, user.Username)
	return botMessage(channel, fmt.Sprintf("%s left %s", user, channel))
}

func (channel *Channel) appendMessage(context *Context, message *Message) {
	channel.Mux.Lock()
	defer channel.Mux.Unlock()
	channel.Messages = append(channel.Messages, message)
	context.Logger.Debugf("[%s] %s", channel, message)
}

func (channel *Channel) userJoined(user *User) *Message {
	channel.Mux.Lock()
	defer channel.Mux.Unlock()
	channel.Users[user.Username] = user
	return botMessage(channel, fmt.Sprintf("%s joined %s", user, channel))
}

func (channel *Channel) String() string {
	return fmt.Sprintf("#%s", channel.Name)
}
