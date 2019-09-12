package utils

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/wonderivan/logger"
	"io"
	"os"
)

var cli *client.Client

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
}
func DockerPull(imageName string) {
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)
	if IsHubDocker(imageName) {
		imageName = "docker.io/library/" + imageName
	}
	reader, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	_, _ = io.Copy(os.Stdout, reader)
	if err != nil {
		logger.Emer("[sealbuild]docker pull image : %s ;error : %s", imageName, err)
	} else {
		logger.Info("[sealbuild]docker pull image : %s success.", imageName)
	}
}
