package api

import (
	v1types "k8s.io/client-go/1.4/pkg/api/v1"
)

// Config is a map of key/value strings which specify the environment variables that should exist
// in the execution environment.
type Config struct {
	App    *App             `json:"-"`
	Values []v1types.EnvVar `json:"values"`
}
