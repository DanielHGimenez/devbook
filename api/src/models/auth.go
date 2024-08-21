package models

import (
	"errors"
	"strings"
)

type Authentication struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (authentication *Authentication) Validar() error {
	var builder strings.Builder

	if authentication.Email == "" {
		builder.WriteString("email shouldn't be null. ")
	}

	if authentication.Password == "" {
		builder.WriteString("password shouldn't be null. ")
	}

	if builder.Len() > 0 {
		return errors.New(builder.String())
	}
	return nil
}
