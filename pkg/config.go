package pkg

import (
	"bytes"
	"github.com/wonderivan/logger"
	"os"
	"text/template"
)

const templateText = string(`{ 
   "name": "{{.Name}}",
   "shell" : "{{.Shell}}"
}`)

func writeFile(fileName string, data []byte) {
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if file != nil {
		_, _ = file.Write(data)
	}
}

func templateContent(name, shell string) []byte {
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
	envMap["Name"] = name
	envMap["Shell"] = shell
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return buffer.Bytes()
}
