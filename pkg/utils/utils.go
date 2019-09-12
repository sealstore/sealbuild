package utils

import "strings"

const oneMBByte = 1024 * 1024

func IsHubDocker(imageName string) string {
	hubName := strings.Split(imageName, "/")
	if strings.Contains(hubName[0], ".") {
		return imageName
	}
	return "docker.io/library/" + imageName
}
