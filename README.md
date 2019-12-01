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
kustomize=false
manifests=/home/cuisongliu/kust/tekton/manifests
images=conf/images.json
files=
```

配置项 | 描述 
:---|:---
config.path| 构建后的包存放的位置
app.enable| app构建是否开启
app.name | 包名称 例如tekton
app.version | 包的版本号，生成的包会带这个版本号 例如1.1 生成的包则为 tekton1.1.tar.gz
app.kustomize | 是否为kustomize脚本
app.manifests | 包的部署yaml文件所在的本地目录，不配置则不会替换模板中的 .Shell 和 .Delete 变量
app.images | app所依赖的所有离线镜像地址，为json格式，格式参照 [images.json](conf/images.json)。不配置则不会替换模板中的.Load .Remove变量。
app.files | 额外需要传入的文件,多个文件以逗号隔开。