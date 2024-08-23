## go打包vue生成的静态资源到exe中
### 安装依赖
```shell
go get -u github.com/jteeuwen/go-bindata/...
go get -u github.com/elazarl/go-bindata-assetfs/...
```

### 执行命令打包静态资源到单独go文件
```
go-bindata-assetfs -o=webapp/webapp.go -pkg=webapp webapp/...
```
> o=webapp/webapp.go为生成的go文件，在webapp目录下的webapp.go
> -pkg=webapp为go文件的包名为webapp
> 打包完成之后会在 webapp 目录下生成一个webapp.go的go文件

### 抓换icon
1. go install github.com/cratonica/2goarray@latest   会安装2goarray.exe 到 gopath/bin
2. .\make_icon.bat .\favicon.ico


### 记录
1. before "github.com/getlantern/systray"   7579KB