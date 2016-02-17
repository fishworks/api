package api

// Build is an executable bundle of software built from source at a specific version or commit
// specified by the deployment process.
//
// To remain compatible with Deis v1, this is just a reference to the fully qualified docker image
// stored in the remote registry. In future iterations, this could be a reference to a rootfs slug
// stored directly in the controller's filesystem or on some other remote store (S3, for example).
type Build struct {
	App *App `json:"-"`
	// Artifact is the fully qualified name of the build along with the version or commit
	// specified by the deployment process. For example, with Docker this would be the
	// fully-qualified docker image name stored on the registry. For rkt, this would be
	// a URL to the Application Container Image (or ACI for short).
	Artifact string `json:"artifact"`
	// Procfile is a process mapping between the artifact's process types and the command
	// associated with said process type.
	Procfile map[string]string `json:"procfile"`
}

func (b *Build) String() string {
	return b.Artifact
}
