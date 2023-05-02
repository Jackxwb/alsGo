package service

import (
	"github.com/gorilla/websocket"
	"golang/service/app"
	"log"
	"strings"
)

// 监听服务 - 打印发送过来的信息
func WebSocketListen(conn *websocket.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[WebSocketListen] Recovered from panic:", err)
			// 处理 panic 错误
			log.Println("出现异常！连接可能已经断开", conn != nil)
		}
	}()

	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		//fmt.Println(string(p))

		log.Printf("收到 %d - %s , %s", messageType, p, strings.Split(string(p), "|"))

		//if err := conn.WriteMessage(messageType, p); err != nil {
		//	log.Println(err)
		//	return
		//}

		inputText := string(p)
		funC := strings.Split(inputText, "|")
		switch funC[0] {
		case "1":
			if len(funC) > 1 {
				go appPing(conn, funC[1])
			} else {
				go appPing(conn, "")
			}
			break
		case "4":
			go appIprtf3Start(conn)
			break
		case "5":
			if len(funC) > 1 {
				go appSpeedTest(conn, funC[1])
			} else {
				go appSpeedTest(conn, "")
			}
			break
		}
	}
}

func appIprtf3Start(conn *websocket.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[appIprtf3Start] Recovered from panic:", err)
			// 处理 panic 错误
		}
	}()
	//log.Println("发起iperf3测试", conn.RemoteAddr().String())
	//log.Println("创建 Iprtf3")
	iperf3 := app.App_Iperf3{}
	//log.Println("初始化 Iprtf3")
	err := iperf3.InitIperf3(conn)
	if err != nil {
		log.Println("Iprtf3 初始化失败 ", err)
		return
	}
	log.Println("Iprtf3 测试", conn.RemoteAddr().String(), iperf3.TaskId)
	//执行方法
	iperf3.Run()
}

func appSpeedTest(conn *websocket.Conn, serverId string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered from panic:", err)
			// 处理 panic 错误
		}
	}()

	speedTest := app.App_SpeedTest{}
	speedTest.Init_SpeedTest(conn, serverId)

	speedTest.JoinTestList()
}

func appPing(conn *websocket.Conn, tager string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered from panic:", err)
			// 处理 panic 错误
		}
	}()

	ping := app.App_Ping{}
	ping.InitPing(conn, tager)

	ping.RunPing()
}
