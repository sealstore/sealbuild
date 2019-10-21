package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/sealstore/sealbuild/pkg/utils"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"os"
	"strings"
)

func app() {
	config := utils.VarsConfig
	//生成的文件名称
	fileName := config.AppName + config.AppVersion + ".tar.gz"
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
	//config.json
	var shell string
	if utils.VarsConfig.AppKustomize {
		shell = "kubectl apply -k manifests"
	} else {
		shell = "kubectl apply -f manifests"
	}
	writeFile(tmpAppDirName+"/config", templateContent(shell, strings.Join(images, " ")))
	//manifests
	_ = os.Mkdir(tmpAppDirName+"/manifests", 0755)
	_ = utils.CopyDir(config.AppManifests, tmpAppDirName+"/manifests")
	//tar
	tmpAppDir, _ := os.Open(tmpAppDirName)
	var tarFiles []*os.File
	tarFiles = append(tarFiles, tmpAppDir)
	logger.Info("[globals]开始创建压缩包。")
	err = utils.Compress(tarFiles, config.Path+"/"+fileName)
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
