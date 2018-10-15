package main

import "fmt"

// User holds user properties
type User struct {
	Username string
}

func (user *User) String() string {
	return fmt.Sprintf("@%s", user.Username)
}
