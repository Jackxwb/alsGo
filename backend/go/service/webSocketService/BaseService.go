package webSocketService

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"golang/Config"
	"golang/service"
	"log"
)

func SendServerBase(conn *websocket.Conn, ips []string) {
	buff := Config.Config.BaseInfo

	remoteAddr := conn.RemoteAddr()
	buff.ClientIp = remoteAddr.String()

	if len(ips) > 1 {
		buff.ClientIp = fmt.Sprintf("%s", ips)
	}

	bytes, err := json.Marshal(buff)
	if err != nil {
		log.Println("生成系统基本信息失败！")
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("1000|"+string(bytes)))

	if err != nil {
		log.Printf("发送消息失败 - SendServerBase - %s", err)
	}
}

func SendHistoricalTraffic(conn *websocket.Conn) {
	frame := service.GetNetDataFrame()
	bytes, err := json.Marshal(frame)
	if err != nil {
		log.Println("获取流量历史帧失败！")
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("101|"+string(bytes)))
	if err != nil {
		log.Printf("发送消息失败 - SendHistoricalTraffic - %s", err)
	}
}
