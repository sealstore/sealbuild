package utils

var VarsConfig *Config

type Config struct {
	Path string `key:"config.path"`
	//app
	AppEnable    bool   `key:"app.enable"`
	AppName      string `key:"app.name"`
	AppVersion   string `key:"app.version"`
	AppKustomize bool   `key:"app.kustomize",default:"false"`
	AppManifests string `key:"app.manifests"`
	AppImages    string `key:"app.images"`
	AppFiles     string `key:"app.files"`
}
