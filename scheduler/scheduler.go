package scheduler

import (
	"errors"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

type State int

type Scheduler interface {
	Create(name, artifact string, command *exec.Cmd) error
	Destroy(name string) error
	Run(name string, command *exec.Cmd) error
	Start(name string) error
	State(name string) State
	Stop(name string) error
}

const (
	StatePending State = iota
	StateRunning
	StateSucceeded
	StateFailed
	StateUnknown
)

func New(module string) (Scheduler, error) {
	log.Debugf("creating scheduler with module %s", module)
	switch module {
	case "docker":
		return NewDockerScheduler()
	case "mock":
		return &MockScheduler{}, nil
	default:
		return nil, errors.New("no scheduler found for type " + module)
	}
}
