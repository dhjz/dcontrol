## 手机控制电脑软件
### 开发语言
- 前端: vue3
- 后端: go
- 使用go:embed 打包前端资源到exe中, 避免分离运行
- 按键传输采用`http`请求, 鼠标移动采用`websocket`
```shell
go get github.com/getlantern/systray
go get github.com/spf13/viper
go get github.com/gorilla/websocket
go get github.com/kbinani/screenshot
go mod tidy
```
### 功能介绍
- go语言控制鼠标移动, go语言控制键盘输入
- 手机控制电脑, 包括按键, 组合键, 音量媒体控制、简易鼠标、简易26键键盘等
- web浏览器端的形式，微信里面打开局域网页面就行
- 支持手机扫描电脑端二维码直接访问局域网地址
- 支持配置自定义本地应用, 手机端一键打开
- 支持简易查看远程桌面
```yml
# congfig.yml
port: 666
apps:
  - name: 微信
    path: E:\Program Files (x86)\Tencent\WeChat\WeChat.exe
  - name: 网易云
    path: E:\Program Files (x86)\NetEase\CloudMusic\cloudmusic.exe
```

### 打包安装
- 双击build.bat 可以打包window运行文件, 双击生成的dist/dcontrol.exe, 即可运行
- 打开页面`http://localhost:666/` 访问页面
- dcontrol.exe -p 666 可以指定端口
### 其他说明
- 页面效果图见`appimg`目录
- ![alt 页面控制图](https://gcore.jsdelivr.net/gh/dhjz/dcontrol@master/appimg/app.jpg)
- ![alt 页面控制图1](https://gcore.jsdelivr.net/gh/dhjz/dcontrol@master/appimg/app1.jpg)
- ![alt 页面控制图2](https://gcore.jsdelivr.net/gh/dhjz/dcontrol@master/appimg/app2.jpg)
