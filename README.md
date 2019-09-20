# sealbuild
build offline packages

# build cloud kernel or app packages
```
sealos build --raw https://dl.k8s.io/v1.15.3/kubernetes-server-linux-amd64.tar.gz --kubeadm-url https://github.com/fanux/kube/releases/download/v1.15.2-lvscare/kubeadm
```



## app build
以下为app构建的主要配置文件
```ini
[config]
path=/home/cuisongliu/app
[oss]
endpoint=oss-cn-beijing.aliyuncs.com
akId=5K52I87TzNPNav8Y
akSk=OrC3gL7dpbC6ZX1DOMUWMczaYVyGYO
bucket=cuisongliu
[app]
enable=true
name=tekton
version=
manifests=/home/cuisongliu/kust/tekton/manifests
images=conf/images.json
ossEnable=true
ossObject=tekton/
```

配置项 | 描述 
:---|:---
config.path| 构建后的包存放的位置
oss.endpoint| 上传包到对应oss的上传地址
oss.akId| 阿里云OSS的accessKey id
oss.akSk| 阿里云OSS的accessKey secure
oss.bucket| 阿里云OSS创建的桶地址
app.enable| app构建是否开启
app.name | 包名称 例如tekton
app.version | 包的版本号，生成的包会带这个版本号 例如1.1 生成的包则为 tekton1.1.tar.gz
app.kustomize | 是否为kustomize脚本
app.manifests | 包的部署yaml文件所在的本地目录
app.images | app所依赖的所有离线镜像地址，为json格式，格式参照 [images.json](conf/images.json)
app.ossEnable | 是否oss上传
app.ossObject | oss上传的目录

