package pkg

import (
	"bytes"
	"github.com/wonderivan/logger"
	"os"
	"text/template"
)

const templateText = string(`LOAD tar -zxvf images.tar.gz && docker load -i images.tar
APPLY {{.Shell}}
DELETE kubectl delete -f manifests
REMOVE sleep 10 && docker rmi -f {{.Images}}`)

func writeFile(fileName string, data []byte) {
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if file != nil {
		_, _ = file.Write(data)
	}
}

func templateContent(shell, images string) []byte {
	tmpl, err := template.New("text").Parse(templateText)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("模板转换错误:", err)
		}
	}()
	if err != nil {
		panic(1)
	}
	var envMap = make(map[string]interface{})
	envMap["Images"] = images
	envMap["Shell"] = shell
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return buffer.Bytes()
}
