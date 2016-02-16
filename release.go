package api

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/fishworks/api/scheduler"
	"github.com/fishworks/api/settings"
)

// Release represents a snapshot of an application's build and config artifacts, which is
// immediately ready for execution in the execution environment.
//
// Releases are an append-only ledger and a release cannot be mutated once it is created.
// Any change must create a new release.
type Release struct {
	App     *App    `json:"-"`
	Build   *Build  `json:"-"`
	Config  *Config `json:"-"`
	Version int     `json:"version"`
}

func (r *Release) String() string {
	return fmt.Sprintf("%s_v%d", r.App.ID, r.Version)
}

// Publish publishes the release to the scheduler.
func (r *Release) Publish() error {
	if r.Build == nil {
		return errors.New("cannot publish; no build associated with this release")
	}
	sched, err := scheduler.New(settings.Scheduler)
	if err != nil {
		return fmt.Errorf("could not publish release: %v", err)
	}
	for typ, command := range r.Build.Procfile {
		id := fmt.Sprintf("%s.%s.1", r.String(), typ)
		if err := sched.Create(
			id,
			r.Build.Artifact,
			exec.Command("sh", "-c", command)); err != nil {
			return err
		}
		if err := sched.Start(id); err != nil {
			return err
		}
		if sched.State(id) != scheduler.StateRunning {
			return errors.New(fmt.Sprintf("job ID %s is flapping", id))
		}
	}
	return nil
}
