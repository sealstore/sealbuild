package utils

var VarsConfig *Config

type Config struct {
	Path      string `key:"config.path"`
	DockerAPI string `key:"config.dockerApi",default:"v1.37"`
	//oss
	OssEndpoint string `key:"oss.endpoint"`
	OssAkId     string `key:"oss.akId"`
	OssAkSk     string `key:"oss.akSk"`
	OssBucket   string `key:"oss.bucket"`
	//app
	AppEnable    bool   `key:"app.enable"`
	AppName      string `key:"app.name"`
	AppVersion   string `key:"app.version"`
	AppManifests string `key:"app.manifests"`
	AppImages    string `key:"app.images"`
	AppOssEnable bool   `key:"app.ossEnable"`
	AppOssObject string `key:"app.ossObject"`
	//cloudkernel
	CloudKernelEnable    bool   `key:"cloudkernel.enable"`
	CloudKernelVersion   string `key:"cloudkernel.version"`
	CloudKernelRaw       string `key:"cloudkernel.raw"`
	CloudKernelKubeadm   string `key:"cloudkernel.kubeadm"`
	CloudKernelOssEnable bool   `key:"cloudkernel.ossEnable"`
	CloudKernelOssObject string `key:"cloudkernel.ossObject"`
}
