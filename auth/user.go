package auth

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user within the Deis authentication system. Username, password and email are
// required. Other fields are optional.
type User struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	// A hash of the password (Deis does not store raw passwords). Raw passwords can be
	// arbitrarily long and can contain any character.
	Password    []byte
	Groups      []string
	Permissions []string
	// Designates that this user has all permissions without explicitly assigning them.
	IsSuper    bool
	IsActive   bool
	LastLogin  time.Time
	DateJoined time.Time
}

// NewUser creates a new User with the specified username, email and password. The password is
// saved as a hash, rather than the raw password.
func NewUser(username, email string, password []byte) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not generate password: %v", err)
	}
	return &User{
		Username:   username,
		Email:      email,
		Password:   hashedPassword,
		IsActive:   true,
		DateJoined: time.Now(),
	}, nil
}

// CheckPassword determines if the supplied password is correct.
func (u User) CheckPassword(password []byte) bool {
	if err := bcrypt.CompareHashAndPassword(u.Password, password); err != nil {
		return false
	}
	return true
}
