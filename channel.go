package main

import "time"

// Channel holds properties related to channels
type Channel struct {
	Name string
}

// Messages are stored as a slice in a map keyed by the channel name
var channelsMessages = map[*Channel][]*Message{}

func (channel *Channel) appendMessage(message string, user *User) (msg *Message) {
	msg = &Message{
		Channel: channel,
		Message: message,
		Time:    time.Now().UTC(),
		User:    user,
	}
	channelsMessages[channel] = append(channelsMessages[channel], msg)
	return
}

func (channel *Channel) String() string {
	return channel.Name
}
