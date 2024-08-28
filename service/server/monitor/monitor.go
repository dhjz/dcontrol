package monitor

import (
	// pointController "express/server/business/point/controller"

	"dcontrol/server/base"
	"dcontrol/server/keys"
	"dcontrol/server/setting"
	"dcontrol/server/utils"
	"fmt"
	"net/http"
	"strings"
)

func HandleApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("Content-Type", "application/json")

	// 获取请求路径   strings.HasSuffix
	path := r.URL.Path
	fmt.Println("redis HandleApi path:", path)

	switch {
	case strings.Contains(path, "/getKeyMap"):
		getKeyMap(w, r)
	case strings.Contains(path, "/getIp"):
		getIp(w, r)
	case strings.Contains(path, "/getApps"):
		getApps(w, r)
	case strings.Contains(path, "/sendkey"):
		sendkey(w, r)
	case strings.Contains(path, "/sendtext"):
		sendtext(w, r)
	case strings.Contains(path, "/open"):
		open(w, r)
	default:
		http.NotFound(w, r)
	}
}

func getKeyMap(w http.ResponseWriter, r *http.Request) {
	// port, err := strconv.Atoi(r.URL.Query().Get("port"))
	base.R(w).Ok(keys.KeyMap)
}

func getIp(w http.ResponseWriter, r *http.Request) {
	base.R(w).Ok(utils.GetMainIP())
}

func getApps(w http.ResponseWriter, r *http.Request) {
	base.R(w).Ok(setting.Conf.Apps)
}

func sendkey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "RBUTTON" {
		keys.ClickMouse("R")
		return
	} else if key == "LBUTTON" {
		keys.ClickMouse("L")
		return
	} else if key == "MBUTTON" {
		keys.ClickMouse("M")
		return
	}
	keys.Run(key)

	base.R(w).Ok(key)
}

func sendtext(w http.ResponseWriter, r *http.Request) {
	val := r.URL.Query().Get("val")
	err := keys.WriteAll(val)
	if err != nil {
		return
	}
	keys.RunKeys(keys.KeyMap["CTRL"], keys.KeyMap["V"])

	base.R(w).Ok(val)
}

func open(w http.ResponseWriter, r *http.Request) {
	cmd1 := r.URL.Query().Get("cmd1")
	cmd2 := r.URL.Query().Get("cmd2")
	utils.RunCmd(cmd1, cmd2)

	base.R(w).Ok(cmd1 + "," + cmd2)
}
