package pkg

import (
	"bytes"
	"github.com/wonderivan/logger"
	"os"
	"strings"
	"text/template"
)

const templateText = string(`LOAD {{.Load}}
APPLY {{.Shell}}
DELETE {{.Delete}}
REMOVE {{.Remove}}`)

func TemplateText() string {
	var sb strings.Builder
	sb.Write([]byte(templateText))
	return sb.String()
}

func writeFile(fileName string, data []byte) {
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if file != nil {
		_, _ = file.Write(data)
	}
}

func templateContent(templateContent, shell, images string) []byte {
	tmpl, err := template.New("text").Parse(templateContent)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("模板转换错误:", err)
		}
	}()
	if err != nil {
		panic(1)
	}
	var envMap = make(map[string]interface{})
	envMap["Remove"] = "sleep 10 && docker rmi -f " + images
	envMap["Shell"] = shell
	envMap["Load"] = "tar -zxvf images.tar.gz && docker load -i images.tar"
	envMap["Delete"] = "kubectl delete -f manifests"
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return buffer.Bytes()
}
