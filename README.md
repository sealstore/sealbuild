# sealbuild
build offline packages



## app build
以下为app构建的主要配置文件
```ini
[config]
path=/home/cuisongliu/app
[app]
enable=true
name=tekton
version=
manifests=/home/cuisongliu/kust/tekton/manifests
images=conf/images.json
```

配置项 | 描述 
:---|:---
config.path| 构建后的包存放的位置
app.enable| app构建是否开启
app.name | 包名称 例如tekton
app.version | 包的版本号，生成的包会带这个版本号 例如1.1 生成的包则为 tekton1.1.tar.gz
app.manifests | 包的部署yaml文件所在的本地目录
app.images | app所依赖的所有离线镜像地址，为json格式，格式参照 [images.json](conf/images.json)
app.ossEnable | 是否oss上传
app.ossObject | oss上传的目录

