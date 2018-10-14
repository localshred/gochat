package main

// User holds user properties
type User struct {
	Username string
}

func (user *User) String() string {
	return user.Username
}
