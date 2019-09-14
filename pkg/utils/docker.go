package utils

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/wonderivan/logger"
	"io"
	"os"
)

func cli() *client.Client {
	var err error
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return cli
}

func DockerPull(imageName string) {
	cli := cli()
	defer cli.Close()
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)
	reader, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		logger.Emer("[sealbuild]docker pull image : %s ;error : %s", imageName, err)
	} else {
		_, _ = io.Copy(os.Stdout, reader)
		logger.Info("[sealbuild]docker pull image : %s success.", imageName)
	}
}

func DockerRmi(imageName string) {
	cli := cli()
	defer cli.Close()
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)
	items, err := cli.ImageRemove(ctx, imageName, types.ImageRemoveOptions{})
	if err != nil {
		logger.Emer("[sealbuild]docker rmi image : %s ;error : %s", imageName, err)
		return
	}
	for _, item := range items {
		logger.Info("[sealbuild]docker rmi image : %s [Deleted: %s] , [ Untagged : %s] ", imageName, item.Deleted, item.Untagged)
	}
}

func DockerSave(tarName string, images []string) {
	cli := cli()
	defer cli.Close()
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)
	reader, err := cli.ImageSave(ctx, images)
	if err != nil {
		logger.Emer("[sealbuild]docker save image : %s ;error : %s", images, err)
	} else {
		file, err := os.OpenFile(tarName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		defer file.Close()
		defer func() {
			if r := recover(); r != nil {
				logger.Emer("[sealbuild]docker save image : %s ;error : %s", images, err)
			}
		}()
		if err != nil {
			panic(1)
		}
		_, _ = io.Copy(file, reader)
		logger.Info("[sealbuild]docker save image : %s to %s success.", images, tarName)
	}
}

func DockerList() []string {
	cli := cli()
	defer cli.Close()
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)
	images, _ := cli.ImageList(ctx, types.ImageListOptions{})
	var ids []string
	for _, image := range images {
		logger.Info("[sealbuild]docker list image: %s.", image.ID)
		ids = append(ids, image.ID)
	}
	return ids
}

func DockerLoad(tarName string) {
	cli := cli()
	defer cli.Close()
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)
	file, err := os.OpenFile(tarName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	defer file.Close()
	defer func() {
		if r := recover(); r != nil {
			logger.Emer("[sealbuild]docker load image : %s ;error : %s", tarName, err)
		}
	}()
	if err != nil {
		panic(1)
	}
	writer, err := cli.ImageLoad(ctx, file, true)
	if err != nil {
		logger.Emer("[sealbuild]docker load image : %s ;error : %s", tarName, err)
	} else {
		_, _ = io.Copy(file, writer.Body)
		logger.Info("[sealbuild]docker load image : %s success.", tarName)
	}

}

func DockerTag(sourceImageName, targetImageName string) {
	cli := cli()
	defer cli.Close()
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)
	err := cli.ImageTag(ctx, sourceImageName, targetImageName)
	if err != nil {
		logger.Emer("[sealbuild]docker tag  error : %s", err)
	} else {
		logger.Info("[sealbuild]docker tag  success : [source : %s],[target : %s].", sourceImageName, targetImageName)
	}
}

func DockerLogin(userName, password string) bool {
	cli := cli()
	defer cli.Close()
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)
	body, err := cli.RegistryLogin(ctx, types.AuthConfig{
		Username: userName,
		Password: password,
	})
	if err != nil {
		logger.Emer("[sealbuild]docker login  error : %s", err)
		return false
	} else {
		logger.Info("[sealbuild]docker login result : %s.", body.Status)
		return true
	}
}

func DockerPush(imageName, loginToken string) {
	cli := cli()
	defer cli.Close()
	ctx := context.Background()
	cli.NegotiateAPIVersion(ctx)
	reader, err := cli.ImagePush(context.Background(), imageName, types.ImagePushOptions{
		RegistryAuth: loginToken,
	})
	if err != nil {
		logger.Emer("[sealbuild]docker push image : %s ;error : %s", imageName, err)
	} else {
		_, _ = io.Copy(os.Stdout, reader)
		logger.Info("[sealbuild]docker push image : %s success.", imageName)
	}
}
