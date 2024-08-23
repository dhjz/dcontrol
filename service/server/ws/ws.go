package ws

import (
	"dcontrol/server/keys"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// 参考 https://github.com/gorilla/websocket/blob/main/examples/command/main.go

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	// pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	// closeGracePeriod = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	}}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("Client connected from IP: %s\n", r.RemoteAddr)
	// fmt.Printf("User-Agent: %s\n Origin: %s\n", r.Header["User-Agent"], r.Header["Origin"])

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	fmt.Println("ServeWs connected......")

	defer ws.Close()

	// ws.SetReadLimit(maxMessageSize)
	// ws.SetReadDeadline(time.Now().Add(pongWait))
	// ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		// 读取消息
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}

		// 打印接收到的消息
		fmt.Printf("ws Received: %s\n", msg)
		wsdata := string(msg)
		if wsdata == "pos,click" {
			// go keys.RunKeys(keys.KeyMap["LBUTTON"])
			keys.ClickMouse("L")
		} else if wsdata == "pos,longclick" {
			keys.ClickMouse("R")
		} else if strings.HasPrefix(wsdata, "pos,start") {
			parts := strings.Split(wsdata, ",")
			if len(parts) == 4 {
				fx, _ := strconv.ParseFloat(parts[2], 64)
				fy, _ := strconv.ParseFloat(parts[3], 64)
				keys.SetMouse(int(fx), int(fy), true)
			}
		}

		// 可以选择回送消息给客户端
		err = ws.WriteMessage(messageType, msg)
		if err != nil {
			fmt.Println("Error while writing message:", err)
			break
		}
	}
}
