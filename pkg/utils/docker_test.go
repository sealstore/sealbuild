package utils

import (
	"testing"
)

func TestDockerPull(t *testing.T) {
	DockerPull("nginx:alpine")
	//DockerRmi("nginx:alpine")
}
