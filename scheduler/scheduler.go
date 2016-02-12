package scheduler

import (
	"errors"
	"os/exec"
)

const (
	StatePending State = iota
	StateRunning
	StateSucceeded
	StateFailed
	StateUnknown
)

type Scheduler interface {
	Create(name string) error
	Destroy(name string) error
	Run(name string, command *exec.Cmd) error
	Start(name string) error
	State(name string) (State, error)
	Stop(name string) error
}

type State int

func New(module string) (Scheduler, error) {
	switch(module) {
	case "mock":
		return &MockScheduler{}, nil
	default:
		return nil, errors.New("no scheduler found for type " + module)
	}
}
