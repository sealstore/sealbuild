package pkg

import (
	"encoding/json"
	"github.com/fanux/sealbuild/pkg/utils"
	"github.com/wonderivan/logger"
	"io/ioutil"
)

func app() {
	//生成的文件名称
	fileName := utils.VarsConfig.AppName + utils.VarsConfig.AppVersion + ".tar.gz"
	logger.Warn("app:%s", fileName)
	imagesData, _ := ioutil.ReadFile(utils.VarsConfig.AppImages)
	logger.Debug("\njson:%s", string(imagesData))
	images := []string{}
	_ = json.Unmarshal(imagesData, images)
	//for _,v:=range images.image {
	//	logger.Info("images:%s",v)
	//}
}
