package auth

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
)

// Token is an authentication strategy used to authenticate users.
type Token struct {
	User    *User
	Key     string
	Created time.Time
}

// NewToken generates a new auth token for a specified user.
func NewToken(user *User) *Token {
	return &Token{
		User:    user,
		Key:     generateKey(),
		Created: time.Now(),
	}
}

func generateKey() string {
	return uuid.New()
}
