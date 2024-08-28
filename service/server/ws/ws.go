package ws

import (
	"dcontrol/server/keys"
	"dcontrol/server/utils"
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

var startX int
var startY int

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
		_, msg, err := ws.ReadMessage()
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
			arr := strings.Split(wsdata, ",")
			if len(arr) == 4 {
				currX, _ := strconv.Atoi(arr[2])
				currY, _ := strconv.Atoi(arr[3])
				if startX != 0 && (startX != currX || startY != currY) {
					keys.SetMouse(int(currX-startX), int(currY-startY), true)
				}
				startX = currX
				startY = currY
			}
		} else if wsdata == "pos,end" {
			startX = 0
			startY = 0
		}

		// 可以选择回送消息给客户端
		// err = ws.WriteMessage(messageType, msg)
		// if err != nil {
		// 	fmt.Println("Error while writing message:", err)
		// 	break
		// }
		if strings.HasPrefix(wsdata, "screen") {
			arr := strings.Split(wsdata, ",")
			var quality = 75
			if len(arr) == 2 {
				q, _ := strconv.Atoi(strings.TrimSpace(arr[1]))
				quality = q
			}
			imgData, err := utils.CaptureScreen(quality)
			if err != nil {
				fmt.Println("Error capturing screen:", err)
				continue
			}

			err = ws.WriteMessage(websocket.BinaryMessage, imgData)
			if err != nil {
				fmt.Println("Error sending image:", err)
				return
			}
		}
	}
}
