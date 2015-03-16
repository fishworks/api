package api

import (
	"fmt"
	"os"
)

// Build is an executable bundle of software built from source at a specific version or commit
// specified by the deployment process.
//
// To remain compatible with Deis v1, this is just a reference to the fully qualified docker image
// stored in the remote registry. In future iterations, this could be a reference to a rootfs slug
// stored directly in the controller's filesystem or on some other remote store (S3, for example).
type Build struct {
	// Name is the fully qualified name of the build along with the version or commit specified
	// by the deployment process. In docker terms, this would be the fully-qualified docker
	// image stored on the registry (without the remote registry URL, as this is stored in
	// $REGISTRY_URL).
	Name string
}

func (b Build) String() string {
	return fmt.Sprintf("%s/%s", os.Getenv("REGISTRY_URL"), b.Name)
}
