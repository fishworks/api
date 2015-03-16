package token

import (
	"fmt"
	"net/http"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/fishworks/api/auth"
)

// Token is an authentication strategy used to authenticate users.
type Token struct {
	User    *auth.User
	Key     string
	Created time.Time
}

// New generates a new auth token for a specified user.
func New(user *auth.User) *Token {
	return &Token{
		User:    user,
		Key:     generateKey(),
		Created: time.Now(),
	}
}

// IsAuthenticated determines if the auth token in the request is valid.
func (t Token) IsAuthenticated(r *http.Request) bool {
	if fmt.Sprintf("token %s", t.Key) == r.Header.Get("HTTP_AUTHORIZATION") {
		return true
	}
	return false
}

func generateKey() string {
	return uuid.New()
}
