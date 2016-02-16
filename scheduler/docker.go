package scheduler

import (
	"errors"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
)

type DockerScheduler struct {
	Client *docker.Client
	// Containers maps job IDs to Docker container IDs.
	Containers map[string]string
}

func parseRepositoryURL(repo string) docker.PullImageOptions {
	var registry string
	repo, tag := docker.ParseRepositoryTag(repo)
	if strings.Count(repo, "/") > 1 {
		splitRepo := strings.SplitN(repo, "/", 2)
		registry, repo = splitRepo[0], splitRepo[1]
	}
	return docker.PullImageOptions{
		Repository: repo,
		Registry:   registry,
		Tag:        tag,
	}
}

func NewDockerScheduler() (*DockerScheduler, error) {
	client, err := docker.NewClientFromEnv()
	return &DockerScheduler{
		Client:     client,
		Containers: make(map[string]string),
	}, err
}

func (s *DockerScheduler) Create(name, artifact string, command *exec.Cmd) error {
	log.Debugf("pulling image %s", artifact)
	if err := s.Client.PullImage(parseRepositoryURL(artifact), docker.AuthConfiguration{}); err != nil {
		return err
	}
	log.Debugf("creating container with ID %s", name)
	container, err := s.Client.CreateContainer(
		docker.CreateContainerOptions{
			Name: name,
			Config: &docker.Config{
				Image: artifact,
			},
		})
	if err == nil {
		s.Containers[name] = container.ID
	}
	return err
}

func (s *DockerScheduler) Destroy(name string) error {
	return nil
}

func (s *DockerScheduler) Run(name string, command *exec.Cmd) error {
	return nil
}

func (s *DockerScheduler) Start(name string) error {
	if id, ok := s.Containers[name]; ok {
		log.Debugf("starting container with ID %s", id)
		return s.Client.StartContainer(id, nil)
	} else {
		return errors.New("job ID does not exist")
	}
}

func (s *DockerScheduler) State(name string) State {
	if id, ok := s.Containers[name]; ok {
		container, err := s.Client.InspectContainer(id)
		if err != nil {
			log.Error(err)
			return StateUnknown
		}
		if container.State.Running {
			return StateRunning
		} else if container.State.Paused {
			return StatePending
		} else if container.State.Restarting {
			return StatePending
		} else if !container.State.Running && container.State.ExitCode == 0 {
			return StateSucceeded
		} else if !container.State.Running && container.State.ExitCode != 0 {
			return StateFailed
		} else {
			return StateUnknown
		}
	} else {
		log.Errorf("job ID %s does not exist", name)
		return StateUnknown
	}
}

func (s *DockerScheduler) Stop(name string) error {
	return nil
}
