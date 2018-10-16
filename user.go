package main

import "fmt"

// User holds user properties
type User struct {
	Username string `json:"username"`
}

var botUser = &User{
	Username: "gochatbot",
}

func (user *User) String() string {
	return fmt.Sprintf("@%s", user.Username)
}
