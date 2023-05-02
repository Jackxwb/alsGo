package app

import (
	"bufio"
	"errors"
	"github.com/bytedance/sonic/utf8"
	"github.com/gorilla/websocket"
	"golang/Config"
	"golang/util"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

type BaseRun struct {
	Exe string
	Arg []string
	//超时时间
	TimeOut int
	TaskId  string

	Ws              *websocket.Conn
	SendTTextPrefix string

	DoText func(txt string)
}

func (b *BaseRun) Init() {
	b.TaskId = util.OnlyKey()
	//b.TimeOut = 60
	b.TimeOut = Config.Config.App.DefTimeOut
	b.DoText = func(txt string) {
		b.SendToWs(txt)
	}
}

func (b BaseRun) autoTrTxt(data []byte) ([]byte, error) {
	if !utf8.Validate(data) {
		return util.TextGbkToUtf8(data)
	}
	return data, nil
}

func (b BaseRun) RunExternalPrograms() {
	var command *exec.Cmd
	// 超时自动结束进程
	var timer *time.Timer = nil
	defer func() {
		if err := recover(); err != nil {
			log.Println("[RunExternalPrograms] Recovered from panic:", err)
			// 处理 panic 错误
		}
		//停用定时结束
		if timer != nil {
			defer func() {
				if err := recover(); err != nil {
					//log.Println("Recovered from panic:", err)
				}
			}()
			timer.Stop()
		}
		if command == nil || command.Process == nil {
			log.Println("RunExternalPrograms 运行结束")
		} else {
			log.Println("RunExternalPrograms 运行结束,pid->", command.Process.Pid)
		}
	}()

	command = exec.Command(b.Exe, b.Arg...)

	dir, err := os.Getwd()
	if err == nil {
		command.Dir = dir
	}

	pipe, err := command.StdoutPipe()
	//pipe, err := command.StderrPipe()
	if err != nil {
		log.Println("获取管道失败", err)
		b.SendToWs("获取管道失败")
		return
	}

	//loogTime := false
	//// 超时自动结束进程
	//var timer *time.Timer = nil
	//是否启用
	if b.TimeOut > 0 {
		timer = time.AfterFunc(time.Duration(b.TimeOut)*time.Second, func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println("[time.AfterFunc Auto Close]Recovered from panic:", err)
					// 处理 panic 错误
				}
			}()
			// 进程运行超过3分钟，终止进程
			if command == nil || command.Process == nil || command.Process.Pid == 0 {
				log.Println("运行超时定时器似乎未能被销毁", b.TaskId)
				return
			}
			err := command.Process.Kill()
			if err != nil {
				log.Println("终止进程错误：", err)
				b.SendToWs("超时终止进程错误：" + err.Error())
			}
			log.Println("进程运行超时已被终止", b.Exe, b.TaskId)
			b.SendToWs("进程运行超时已被终止")
			//loogTime = true
		})
	}

	stderrPipe, err := command.StderrPipe()
	if err != nil {
		log.Println("获取Err管道异常")
	} else {
		readErr := func(errPipe io.ReadCloser) {
			defer func() {
				if err := recover(); err != nil {
					log.Println("[readErr]Recovered from panic:", err)
					// 处理 panic 错误
				}
			}()
			reader := bufio.NewReader(errPipe)
			for true {
				line, prefix, err := reader.ReadLine()
				if prefix {
					break
				}
				if err != nil {
					if err == io.EOF {
						//正常结束
						break
					}
					log.Println("管道读取失败", err)
					return
				}
				log.Println("runErr[", command.Process.Pid, "]->", string(line))
			}
		}
		go readErr(stderrPipe)
	}

	log.Println("RunExternalPrograms 运行开始 ", command.Dir, b.Exe, b.Arg)
	if err = command.Start(); err != nil {
		log.Println("运行", b.Exe, "失败,", err)
		b.SendToWs("运行" + b.Exe + "失败," + err.Error())
		return
	}

	//log.Println("EXE->", b.Exe, " Arg->", b.Arg)
	log.Println("pid-> ", command.Process.Pid, "TimeOut->", b.TimeOut)

	reader := bufio.NewReader(pipe)
	for true {
		line, prefix, err := reader.ReadLine()
		if prefix {
			break
		}
		if err != nil {
			if err == io.EOF {
				//正常结束
				break
			}
			log.Println("管道读取失败", err)
			return
		}
		//fmt.Println(line)
		txt, err := b.autoTrTxt(line)
		if err != nil {
			log.Println("转换字符串格式异常", err, "\r\n", line)
			b.DoText(string(line))
			continue
		}
		b.DoText(string(txt))
	}

	if err = command.Wait(); err != nil {
		log.Println("运行", b.Exe, "失败,", err)
		b.SendToWs("运行" + b.Exe + "失败," + err.Error())
		log.Println("EXE->", b.Exe, " Arg->", b.Arg)
		return
	}
	log.Println("运行", b.Exe, "结束", timer == nil) // , timer
}

func (b BaseRun) SendToWs(str string) (err error) {
	defer func() {
		if errR := recover(); errR != nil {
			log.Println("[SendToWs] Recovered from panic:", errR)
			// 处理 panic 错误
			//log.Println("发送数据出现异常")
		}
	}()
	if b.Ws == nil {
		err = errors.New("[SendToWs] webSocket未初始化!")
		return err
	}
	err = b.Ws.WriteMessage(websocket.TextMessage, []byte(b.SendTTextPrefix+str))
	return err
}
func (b BaseRun) SendSourseTextToWs(str string) (err error) {
	defer func() {
		if errR := recover(); errR != nil {
			log.Println("[SendSourseTextToWs] Recovered from panic:", errR)
			// 处理 panic 错误
			//log.Println("[SendSourseTextToWs] 发送数据出现异常")
		}
	}()
	if b.Ws == nil {
		err = errors.New("webSocket未初始化!")
		return err
	}
	err = b.Ws.WriteMessage(websocket.TextMessage, []byte(str))
	return err
}

// 寻找文件是否存在
func (b BaseRun) findFileIsExist(path string) bool {
	_, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Println("判断文件是否存在出现异常", err)
		return false
	}
	return true
}
