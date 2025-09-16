package models

import "fmt"

type JSONMap = map[string]any

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = []User{
	{ID: 1, Username: "root", Password: "root"},
	{ID: 2, Username: "test", Password: "test"},
}

func (u *User) CheckPassword(password string) error {
	if u.Password != password {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func GetUserByUsername(username string) (*User, error) {
	for _, user := range users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("no such user with username %q", username)
}
