package utils

import (
	"github.com/wonderivan/logger"
	"gopkg.in/ini.v1"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type IniFile struct {
	iniFile *ini.File
}

func LoadIni(data interface{}, path string) {
	iniConfigure := loadIniFile(path)
	loadIni(data, iniConfigure)
}

func loadIniFile(path string) *ini.File {
	var err error
	iniFile, err := ini.Load(path)
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]读取配置文件失败,请检测配置文件路径是否存在[" + path + "]")
			os.Exit(1)
		}
	}()
	if err != nil {
		logger.Error(err)
		panic("读取配置文件失败,请检测配置文件路径是否存在[" + path + "]")
	}
	return iniFile
}

func loadIni(data interface{}, iniConfigure *ini.File) {
	tp := reflect.TypeOf(data)
	tv := reflect.ValueOf(data)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	if tv.Kind() == reflect.Ptr {
		tv = tv.Elem()
	}
	fieldNum := tp.NumField()
	for i := 0; i < fieldNum; i++ {
		eid := tp.Field(i)
		vid := tv.Field(i)
		key := eid.Tag.Get("key")
		if key == "" {
			logger.Error("结构体[%s]中的字段%s标签key属性为空,请修改后重试", tp.Name(), eid.Name)
			os.Exit(1)
		}
		keys := priGetKey(key)
		//logger.Alert("读取配置:key:%s", key)
		key = keys[1]
		section := keys[0]
		defaults := eid.Tag.Get("default")
		switch vid.Kind() {
		case reflect.String:
			value := iniConfigure.Section(section).Key(key).MustString(defaults)
			vid.SetString(value)
			logger.Trace("读取配置:[%s]%s=%s", section, key, value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var defaultValue int
			if defaults == "" {
				logger.Trace("结构体[%s]中的字段%s标签defaults属性为空,默认值为0", tp.Name(), eid.Name)
				defaultValue = 0
			} else {
				var err error
				if defaultValue, err = strconv.Atoi(defaults); nil != err {
					defaultValue = 0
				}
			}
			value := iniConfigure.Section(section).Key(key).MustInt(defaultValue)
			vid.SetInt(int64(value))
			logger.Trace("读取配置:[%s]%s=%d", section, key, value)
		case reflect.Bool:
			var defaultValue bool
			if defaults == "" {
				logger.Trace("结构体[%s]中的字段%s标签defaults属性为空,默认值为false", tp.Name(), eid.Name)
				defaultValue = false
			} else {
				var err error
				if defaultValue, err = strconv.ParseBool(defaults); nil != err {
					defaultValue = false
				}
			}
			value := iniConfigure.Section(section).Key(key).MustBool(defaultValue)
			vid.SetBool(value)
			logger.Trace("读取配置:[%s]%s=%t", section, key, value)
		}

	}
}

func priGetKey(key string) []string {
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		keysR := make([]string, 2)
		keysR[0] = ""
		keysR[1] = key
		return keysR
	}
	return keys
}
