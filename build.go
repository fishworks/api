package api

// Build is an executable bundle of software built from source at a specific version or commit
// specified by the deployment process.
//
// To remain compatible with Deis v1, this is just a reference to the fully qualified docker image
// stored in the remote registry. In future iterations, this could be a reference to a rootfs slug
// stored directly in the controller's filesystem or on some other remote store (S3, for example).
type Build struct {
	App *App `json:"-"`
	// Image is the fully qualified name of the image along with the version or commit specified
	// by the deployment process. With Docker, this would be the fully-qualified docker image name.
	Image string `json:"image"`
	// Procfile is a process mapping between the images's process types and the arguments
	// (equivalent to a Docker container's CMD entrypoint) associated with said process type.
	Procfile map[string][]string `json:"procfile"`
}

func (b *Build) String() string {
	return b.Image
}
