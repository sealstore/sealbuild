package utils

import "strings"

const oneMBByte = 1024 * 1024

func IsHubDocker(imageName string) bool {
	hubName := strings.Split(imageName, "/")
	if strings.Contains(hubName[0], ".") {
		return false
	}
	return true
}
