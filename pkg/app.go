package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/fanux/sealbuild/pkg/utils"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"os"
)

func app() {
	//生成的文件名称
	fileName := utils.VarsConfig.AppName + utils.VarsConfig.AppVersion + ".tar.gz"
	logger.Warn("app:%s", fileName)
	imagesData, _ := ioutil.ReadFile(utils.VarsConfig.AppImages)
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
	tmpImageName := fmt.Sprintf("/tmp/images_%s.tar", utils.VarsConfig.AppName+utils.VarsConfig.AppVersion)
	utils.DockerSave(tmpImageName, images)
	tmpAppDirName := fmt.Sprintf("/tmp/%s", utils.VarsConfig.AppName+utils.VarsConfig.AppVersion)
	_ = os.RemoveAll(tmpAppDirName)
	err := os.Mkdir(tmpAppDirName, 0755)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]创建目录失败: /tmp/"+utils.VarsConfig.AppName+utils.VarsConfig.AppVersion, err)
			os.Exit(1)
		}
	}()
	if err != nil {
		panic(1)
	}
	_ = os.Rename(tmpImageName, tmpAppDirName+"/images.tar")
	//config.json
	writeFile(tmpAppDirName+"/config.json", templateContent("app", ""))
	//manifests

	//删除镜像文件
	_ = os.Remove(tmpImageName)
}
