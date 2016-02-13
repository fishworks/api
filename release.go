package api

import (
	"errors"
	"fmt"
)

// Release represents a snapshot of an application's build and config artifacts, which is
// immediately ready for execution in the execution environment.
//
// Releases are an append-only ledger and a release cannot be mutated once it is created.
// Any change must create a new release.
type Release struct {
	Build   *Build  `json:"-"`
	Config  *Config `json:"-"`
	Version int     `json:"version"`
}

func (r *Release) String() string {
	return fmt.Sprintf("v%d", r.Version)
}

// Publish publishes the build along with the config to a remote store.
func (r *Release) Publish() error {
	if r.Build == nil {
		return errors.New("cannot publish; no build associated with this release")
	}
	// TODO: implement scheduler deploy
	return errors.New("not implemented")
}
