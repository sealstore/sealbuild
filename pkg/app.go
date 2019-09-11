package pkg

import (
	"encoding/json"
	"fmt"
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
	var images []string
	//注意数组转换需要加 &
	_ = json.Unmarshal(imagesData, &images)
	for _, v := range images {
		if v != "" {
			logger.Info("images:%s", v)
			cmd := fmt.Sprintf("docker pull %s", v)
			var command string
			var err error
			if command, err = utils.Shell(cmd); err != nil {
				logger.Error(err)
			}
			logger.Alert(command)
		}
	}
}
