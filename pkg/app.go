package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/sealstore/sealbuild/pkg/utils"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func app(templateFile string) {
	config := utils.VarsConfig
	//生成的文件名称
	fileName := config.AppName + config.AppVersion + ".tar"
	logger.Warn("app:%s", fileName)
	imagesData, _ := ioutil.ReadFile(config.AppImages)
	logger.Debug("\njson:%s", string(imagesData))
	var images []string
	//注意数组转换需要加 &
	_ = json.Unmarshal(imagesData, &images)
	for _, v := range images {
		if v != "" {
			utils.DockerPull(v)
		}
	}
	//生成文件
	tmpImageName := fmt.Sprintf("/tmp/images_%s.tar", config.AppName+config.AppVersion)
	utils.DockerSave(tmpImageName, images)
	tmpAppDirName := fmt.Sprintf("/tmp/%s", config.AppName)
	_ = os.RemoveAll(tmpAppDirName)
	err := os.Mkdir(tmpAppDirName, 0755)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]创建目录失败: /tmp/"+config.AppName+config.AppVersion, err)
			os.Exit(1)
		}
	}()
	if err != nil {
		panic(1)
	}
	_ = os.Rename(tmpImageName, tmpAppDirName+"/images.tar")

	tarFile, _ := os.Open(tmpAppDirName + "/images.tar")
	_ = utils.Compress([]*os.File{tarFile}, tmpAppDirName+"/images.tar.gz")
	_ = os.RemoveAll(tmpAppDirName + "/images.tar")

	//config.json
	var shell string
	if utils.VarsConfig.AppKustomize {
		shell = "kubectl apply -k manifests"
	} else {
		shell = "kubectl apply -f manifests"
	}
	var templateFileContent string
	if templateFile == "" {
		templateFileContent = TemplateText()
	} else {
		templateFileData, err := ioutil.ReadFile(templateFile)
		templateFileContent = string(templateFileData)
		defer func() {
			if r := recover(); r != nil {
				logger.Error("[globals]template file read failed:", err)
			}
		}()
		if err != nil {
			panic(1)
		}
	}
	writeFile(tmpAppDirName+"/config", templateContent(templateFileContent, shell, strings.Join(images, " ")))

	//manifests
	_ = os.Mkdir(tmpAppDirName+"/manifests", 0755)
	_ = utils.CopyDir(config.AppManifests, tmpAppDirName+"/manifests")
	//tar
	tmpAppDir, _ := os.Open(tmpAppDirName)
	var tarFiles []*os.File
	tarFiles = append(tarFiles, tmpAppDir)
	logger.Info("[globals]开始创建压缩包。")
	var tarFilesArr []string
	tarFilesArr = append(tarFilesArr, "config")
	tarFilesArr = append(tarFilesArr, "images.tar.gz")
	tarFilesArr = append(tarFilesArr, "manifests")
	shellTar := fmt.Sprintf("cd %s && tar cvf %s %s", tmpAppDirName, fileName, strings.Join(tarFilesArr, " "))
	cmd := exec.Command("/bin/bash", "-c", shellTar)
	err = cmd.Run()
	//err = utils.Tar(tarFilesArr,config.Path+"/"+fileName)
	////err = utils.Compress(tarFiles, config.Path+"/"+fileName)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]创建tar失败: ", err)
			os.Exit(1)
		}
	}()
	if err != nil {
		panic(1)
	}
	_, err = utils.CopyFile(tmpAppDirName+"/"+fileName, config.Path+"/"+fileName)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]创建tar失败: ", err)
			os.Exit(1)
		}
	}()
	if err != nil {
		panic(1)
	}
	logger.Info("[globals]创建压缩包成功。")
}
