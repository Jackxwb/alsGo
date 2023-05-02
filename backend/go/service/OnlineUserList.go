package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

var OnlineList = []*websocket.Conn{}

//type user struct {
//	conn *websocket.Conn
//}

var listLock sync.RWMutex

func AddUser(user *websocket.Conn) {
	defer func() {
		listLock.Unlock()
	}()
	listLock.Lock()

	OnlineList = append(OnlineList, user)
	user.SetCloseHandler(func(code int, text string) error {
		DelUser(user)
		return nil
	})
}

func DelUser(user *websocket.Conn) {
	defer func() {
		listLock.Unlock()
	}()
	listLock.Lock()

	for i := 0; i < len(OnlineList); i++ {
		if OnlineList[i].RemoteAddr().String() == user.RemoteAddr().String() {
			OnlineList = append(OnlineList[:i], OnlineList[i+1:]...)
			log.Println("链接已断开", user.RemoteAddr().String())
			return
		}
	}
	log.Println("未能找到对应链接", user.RemoteAddr().String())
}

func sendText(txt string) {
	defer func() {
		listLock.Unlock()
		if err := recover(); err != nil {
			log.Println("Recovered from panic:", err)
		}
	}()
	listLock.Lock()

	for i := 0; i < len(OnlineList); i++ {
		defer func() {
			if err := recover(); err != nil {
				//log.Println("Recovered from panic:", err)
				fmt.Errorf("Recovered from panic: %s", err)
				// 处理 panic 错误

				//删除当前对象
				OnlineList = append(OnlineList[:i], OnlineList[i+1:]...)
				i--
				fmt.Println(i, len(OnlineList))
			}
		}()
		conn := OnlineList[i]
		//err := conn.WriteMessage(websocket.TextMessage, []byte(txt))
		//if err != nil {
		//	log.Fatalf("发送失败 - sendText - [%s] %s", conn.RemoteAddr(), err)
		//
		//	//尝试断开连接
		//	//util.Try(func() {
		//	//	conn.Close()
		//	//}, func(err interface{}) {})
		//
		//	////删除当前对象
		//	//OnlineList = append(OnlineList[:i], OnlineList[i+1:]...)
		//	//i--
		//	//fmt.Println(i, len(OnlineList))
		//	//
		//	panic(err)
		//}
		writer, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Println("获取 writer 失败！")
			//panic(err)
			closeConnet(conn)
		} else {
			_, err := writer.Write([]byte(txt))
			if err != nil {
				log.Println("writer 写出 失败！")
				//panic(err)
				closeConnet(conn)
			}
			writer.Close()
		}
	}
}

func closeConnet(conn *websocket.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[closeConnet]Recovered from panic:", err)
			// 处理 panic 错误
		}
	}()
	conn.Close()
}

func UpdateNetMessage(data DataFrame) {
	for _, info := range data.Data {
		txt := "100|"
		txt += info.Name
		txt += "|"
		txt += info.Recv
		txt += "|"
		txt += info.Send
		txt += "|"

		sendText(txt)
	}
	//go sendText(txt)
}

func GetOnlineNum() int {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[GetOnlineNum]Recovered from panic:", err)
			// 处理 panic 错误
		}
	}()
	return len(OnlineList)
}
