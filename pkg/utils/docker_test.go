package utils

import (
	"testing"
)

func TestDockerPull(t *testing.T) {
	//DockerPull("nginx")
	//DockerRmi("nginx")
	//DockerSave("/home/cuisongliu/aa.tar",[]string{"nginx"})
	DockerLoad("/home/cuisongliu/aa.tar")
}
