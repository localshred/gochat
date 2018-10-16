package main

import "fmt"

// User holds user properties
type User struct {
	Username string
}

var botUser = &User{
	Username: "bot",
}

func (user *User) String() string {
	return fmt.Sprintf("@%s", user.Username)
}
