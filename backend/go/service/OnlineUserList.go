package service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"golang/util"
	"log"
	"sync"
)

type MyConn struct {
	*websocket.Conn
	Ips     []string
	ErrNum  int
	LastErr string
}

// var OnlineList = []*websocket.Conn{}
var OnlineList = []*MyConn{}

//type user struct {
//	conn *websocket.Conn
//}

type OnlineUser struct {
	Addr    string
	Ips     []string
	ErrNum  int
	LastErr string
}

func GetOnlineUser() []OnlineUser {
	defer func() {
		listLock.RUnlock()
	}()
	listLock.RLock()

	var result []OnlineUser

	for i := 0; i < len(OnlineList); i++ {
		result = append(result, OnlineUser{
			Addr:    OnlineList[i].RemoteAddr().String(),
			Ips:     OnlineList[i].Ips,
			ErrNum:  OnlineList[i].ErrNum,
			LastErr: OnlineList[i].LastErr,
		})
	}

	return result
}

var listLock sync.RWMutex

func AddUser(user *websocket.Conn, ips []string) {
	defer func() {
		listLock.Unlock()
	}()
	listLock.Lock()

	var myUser MyConn
	myUser.Conn = user
	myUser.Ips = ips

	OnlineList = append(OnlineList, &myUser)
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
func DelUserNoLock(user *MyConn) bool {
	for i := 0; i < len(OnlineList); i++ {
		if OnlineList[i].RemoteAddr().String() == user.RemoteAddr().String() {
			OnlineList = append(OnlineList[:i], OnlineList[i+1:]...)
			log.Println("链接已断开", user.RemoteAddr().String())
			return true
		}
	}
	log.Println("未能找到对应链接", user.RemoteAddr().String())
	return false
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
			//conn.LastErr = "获取 writer 失败！"
			//conn.ErrNum++
			//OnlineList[i] = conn
			//
			//log.Println("获取 writer 失败！")
			////panic(err)
			//closeConnet(conn)
			//if DelUserNoLock(conn) {
			//	i--
			//}
			if havErr(conn, "获取 writer 失败！") {
				i--
			} else {
				OnlineList[i] = conn
			}
		} else {
			_, err := writer.Write([]byte(txt))
			if err != nil {
				//log.Println("writer 写出 失败！")
				//closeConnet(conn)
				if havErr(conn, "writer 写出 失败！") {
					i--
				} else {
					OnlineList[i] = conn
				}
			}
			writer.Close()
		}
	}
}
func havErr(conn *MyConn, errText string) bool {
	conn.LastErr = util.GetNowStr() + "获取 writer 失败！"
	conn.ErrNum++

	log.Println("[", conn.ErrNum, "]", conn.LastErr, conn.RemoteAddr().String())
	closeConnet(conn)
	return DelUserNoLock(conn)
}

func closeConnet(conn *MyConn) {
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
