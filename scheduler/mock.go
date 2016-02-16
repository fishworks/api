package scheduler

import (
	"os/exec"
)

type MockScheduler struct{}

func (m *MockScheduler) Create(name, artifact string, command *exec.Cmd) error {
	return nil
}

func (m *MockScheduler) Destroy(name string) error {
	return nil
}

func (m *MockScheduler) Run(name string, command *exec.Cmd) error {
	return nil
}

func (m *MockScheduler) Start(name string) error {
	return nil
}

func (m *MockScheduler) State(name string) State {
	return StateRunning
}

func (m *MockScheduler) Stop(name string) error {
	return nil
}
