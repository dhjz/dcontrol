package main

import (
	"dcontrol/server/base"
	"dcontrol/server/monitor"
	"dcontrol/server/setting"
	"dcontrol/server/utils"
	"dcontrol/server/keys"
	"dcontrol/server/ws"
	"flag"
	"fmt"
	"io/fs"
	"net/http"

	"embed"
)

//go:embed webapp
var f embed.FS

func main() {

	port := flag.Int("p", 0, "server port")
	base.RunPort = *port
	filePath := flag.String("f", "./config.yml", "server config file")
	// dir := flag.String("d", "./webapp", "server static dir")
	flag.Parse()
	//1.加载配置
	setting.Init(*filePath)
	base.RunPort = setting.Conf.Port
	if *port != 0 {
		base.RunPort = *port
	}
	addr := fmt.Sprintf(":%d", base.RunPort)

	http.HandleFunc("/control-api/monitor/", monitor.HandleApi)
	http.HandleFunc("/ws", ws.ServeWs)

	// 注册静态资源
	st, _ := fs.Sub(f, "webapp")
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(st))))

	fmt.Printf("***********************app run on http://localhost:%d/ *******************", base.RunPort)
	fmt.Println("")
	fmt.Println(utils.GetAllIPs())

	if setting.Conf.Open {
		utils.OpenBrowser(fmt.Sprintf("http://localhost:%d/", base.RunPort))
	}
	go func() {
		utils.GenTaskBarIcon()
	}()
	keys.ListenScroll()
	fmt.Println("start http.... ", base.RunPort)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("start http error: ", err)
	}
	fmt.Println("start http success ", base.RunPort)
}
