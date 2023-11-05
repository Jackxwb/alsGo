package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang/entity/api"
	"golang/service"
	"golang/service/webSocketService"
	"golang/util"
	"log"
	"net/http"
)

// websocket 配置
var WsConfig = websocket.Upgrader{
	//缓存 5Mb
	ReadBufferSize:  5 * 1024,
	WriteBufferSize: 5 * 1024,
	//压缩
	EnableCompression: true,
	//跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func DefWebSocket(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered from panic:", err)
			// 处理 panic 错误
		}
	}()

	if !websocket.IsWebSocketUpgrade(c.Request) {
		api.ApiReasultFail("错误的请求类型！").ToContext(c)
		return
	}
	//url := c.Query("url")
	ws, err := WsConfig.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		api.ApiReasultFail(err.Error()).ToContext(c)
		return
	}
	//err = ws.WriteMessage(websocket.TextMessage, []byte("构建开始"))
	//if err != nil {
	//	fmt.Println(err)
	//}

	//获取代理信息
	ipS := util.GetTrueIp(c)

	//发送基本信息
	webSocketService.SendServerBase(ws, ipS)

	//发送缓存
	webSocketService.SendHistoricalTraffic(ws)

	//添加到队列中
	service.AddUser(ws, ipS)

	go service.WebSocketListen(ws)
}
