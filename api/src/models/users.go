package models

import (
	"errors"
	"strings"
	"time"
)

type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (user *User) Validate() error {
	var builder strings.Builder

	if user.Name == "" {
		builder.WriteString("name shouldn't be null. ")
	}

	if user.Nick == "" {
		builder.WriteString("nick shouldn't be null. ")
	}

	if user.Email == "" {
		builder.WriteString("email shouldn't be null. ")
	}

	if user.Password == "" {
		builder.WriteString("password shouldn't be null. ")
	}

	if builder.Len() > 0 {
		return errors.New(builder.String())
	}
	return nil
}
