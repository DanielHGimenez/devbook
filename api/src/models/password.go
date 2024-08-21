package models

import (
	"errors"
	"strings"
)

type PasswordChange struct {
	Current string `json:"current,omitempty"`
	New     string `json:"new,omitempty"`
}

func (p *PasswordChange) Validate() error {
	var builder strings.Builder

	if p.Current == "" {
		builder.WriteString("'current' shouldn't be null. ")
	}

	if p.New == "" {
		builder.WriteString("'new' shouldn't be null. ")
	}

	if builder.Len() > 0 {
		return errors.New(builder.String())
	}
	return nil
}
