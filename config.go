package api

// Config is a map of key/value strings which specify the environment variables that should exist
// in the execution environment. This also includes other values like the maximum allocated memory,
// maximum cpu shares, isolating applications on a set of hosts via tags etc.
type Config struct {
	App         *App              `json:"-"`
	Environment map[string]string `json:"environment"`
}
