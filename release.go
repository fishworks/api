package api

import (
	"fmt"

	"k8s.io/client-go/1.4/kubernetes"
	v1types "k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/rest"
)

var (
	ErrNoBuildToPublish = &ReleaseError{"no build to publish with this release"}
)

type ReleaseError struct {
	Message string
}

func (r *ReleaseError) Error() string {
	return fmt.Sprintf("could not publish release: %s", r.Message)
}

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

// Publish publishes the release to kubernetes.
func (r *Release) Publish() error {
	if r.Build == nil {
		return ErrNoBuildToPublish
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	for typ, command := range r.Build.Procfile {
		podName := fmt.Sprintf("%s_%s", r.String(), typ)
		pod := &v1types.Pod{
			ObjectMeta: v1types.ObjectMeta{
				Name:      podName,
				Namespace: r.App.ID,
				Labels: map[string]string{
					"heritage": "deis",
				},
			},
			Spec: v1types.PodSpec{
				RestartPolicy: v1types.RestartPolicyAlways,
				Containers: []v1types.Container{
					v1types.Container{
						Name:            podName,
						Image:           r.Build.Image,
						ImagePullPolicy: v1types.PullAlways,
						Command:         command,
						Env:             r.Config.Values,
					},
				},
			},
		}
		// Schedule the pod
		if _, err := clientset.Pods(r.App.ID).Create(pod); err != nil {
			return err
		}
	}
	return nil
}
