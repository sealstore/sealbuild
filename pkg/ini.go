package pkg

import (
	"github.com/wonderivan/logger"
	"gopkg.in/ini.v1"
)

type IniFile struct {
	iniFile *ini.File
}
type SubConfigure struct {
	JobName      string
	BuildVersion string
}

type IniConfigure struct {
	User         string
	Password     string
	Host         string
	Artifactory  string
	SubConfigure []*SubConfigure
}

func LoadIniFile(path string) *IniFile {
	var err error
	iniFile, err := ini.Load(path + ".ini")
	defer func() {
		if r := recover(); r != nil {
			logger.Error("[globals]读取配置文件失败,请检测配置文件路径是否存在[" + path + "]")
		}
	}()
	if err != nil {
		logger.Error(err)
		panic("读取配置文件失败,请检测配置文件路径是否存在[" + path + "]")
	}
	return &IniFile{iniFile: iniFile}
}

func (t *IniFile) LoadData() *IniConfigure {
	ss := t.iniFile.SectionStrings()
	//过滤default
	ss = ss[1:]
	iconfig := &IniConfigure{}
	if len(ss) > 1 {
		iconfig.User = t.iniFile.Section("install").Key("username").Value()
		iconfig.Password = t.iniFile.Section("install").Key("password").Value()
		iconfig.Host = t.iniFile.Section("install").Key("host").Value()
		iconfig.Artifactory = t.iniFile.Section("install").Key("artifactory").Value()
	}
	ss = ss[1:]
	var sconfigs []*SubConfigure
	for _, v := range ss {
		sconfig := &SubConfigure{}
		sconfig.JobName = t.iniFile.Section(v).Key("job_name").Value()
		sconfig.BuildVersion = t.iniFile.Section(v).Key("build_version").Value()
		sconfigs = append(sconfigs, sconfig)
	}
	iconfig.SubConfigure = sconfigs
	return iconfig
}
