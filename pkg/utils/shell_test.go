package utils

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"net/http"
	"testing"
)

type FakeRoundTripper struct {
	message  string
	status   int
	header   map[string]string
	requests []*http.Request
}

func TestShell(t *testing.T) {
	//s:=Shell("docker","pull","nginx")
	//t.Log(s)
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	opts := docker.PullImageOptions{
		Repository: "nginx",
		Tag:        "latest",
		Registry:   "https://index.docker.io/v1/",
	}
	err = client.PullImage(opts, docker.AuthConfiguration{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("success")
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		panic(err)
	}
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}

}
