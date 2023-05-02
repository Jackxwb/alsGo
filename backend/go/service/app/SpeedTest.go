package app

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
)

var TestList []*App_SpeedTest
var TestListLock sync.RWMutex

type SpeedTestResult struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Ping      struct {
		Jitter   float64 `json:"jitter"`
		Latency  float64 `json:"latency"`
		Progress float64 `json:"progress"`
	} `json:"ping"`
	Download struct {
		Bandwidth float64 `json:"bandwidth"`
		Bytes     float64 `json:"bytes"`
		Elapsed   float64 `json:"elapsed"`
		Progress  float64 `json:"progress"`
		Latency   struct {
			Iqm float64 `json:"iqm"`
		}
	} `json:"download"`
	Upload struct {
		Bandwidth float64 `json:"bandwidth"`
		Bytes     float64 `json:"bytes"`
		Elapsed   float64 `json:"elapsed"`
		Progress  float64 `json:"progress"`
		Latency   struct {
			Iqm float64 `json:"iqm"`
		}
	} `json:"upload"`
	Isp       string `json:"isp"`
	Interface struct {
		InternalIp string `json:"internalIp"`
		Name       string `json:"name"`
		MacAddr    string `json:"macAddr"`
		IsVpn      bool   `json:"isVpn"`
		ExternalIp string `json:"externalIp"`
	} `json:"interface"`
	Server struct {
		Id       int    `json:"id"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
		Ip       string `json:"ip"`
	} `json:"server"`
	Result struct {
		Id        string `json:"id"`
		Url       string `json:"url"`
		Persisted bool   `json:"persisted"`
	} `json:"result"`
}

type SpeedTestResultType struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
}

var (
	codeStart         = 100
	codeStartPing     = 101
	codeStartDownload = 102
	codeStartUpload   = 103
	codeFinish        = 104
)

type App_SpeedTest struct {
	BaseRun
	IsRun bool
	//前端
	IsStart         bool
	IsStartPing     bool
	IsStartDownload bool
	IsStartUpload   bool
	IsFinish        bool
}

func (s *App_SpeedTest) Init_SpeedTest(conn *websocket.Conn, serverId string) {
	s.Init()
	s.IsRun = false

	s.Exe = "./app/speedtest"
	if !s.findFileIsExist(s.Exe) {
		s.Exe = "./app/speedtest.exe"
		if !s.findFileIsExist(s.Exe) {
			s.Exe = "speedtest"
		}
	}

	s.Ws = conn

	s.SendTTextPrefix = "5|" + s.TaskId + "|"

	s.Arg = []string{"-p", "no", "--accept-license", "-f", "jsonl"}
	if serverId != "" {
		s.Arg = append(s.Arg, "-s", serverId)
	}

	s.DoText = func(txt string) {
		result := SpeedTestResult{}
		err := json.Unmarshal([]byte(txt), &result)
		if err != nil {
			log.Println("解析数据失败![", err, "]", txt)
			return
		}
		switch result.Type {
		case "testStart":
			bytes, err := json.Marshal(result.Server)
			if err != nil {
				log.Println("[testStart]解析数据异常", err, txt)
				break
			}
			s.SendToWs("200|" + string(bytes))
			break
		case "ping":
			if !s.IsStartPing {
				s.IsStartPing = true
				s.sendStartCode(codeStartPing)
			}
			//进度
			s.SendToWs("3|" + strconv.Itoa(int(result.Ping.Progress*100)))
			full := strconv.Itoa(int(result.Ping.Progress * 33))
			s.SendToWs("4|" + full)
			//数据
			s.SendToWs("201|" + strconv.FormatFloat(result.Ping.Latency, 'f', 0, 32))
			break
		case "download":
			if !s.IsStartDownload {
				s.IsStartDownload = true
				s.sendStartCode(codeStartDownload)
			}
			//进度
			s.SendToWs("3|" + strconv.Itoa(int(result.Download.Progress*100)))
			full := strconv.Itoa(33 + int(result.Download.Progress*33))
			s.SendToWs("4|" + full)
			//数据
			s.SendToWs("202|" + strconv.FormatFloat(result.Download.Bandwidth, 'f', 0, 32))
			break
		case "upload":
			if !s.IsStartUpload {
				s.IsStartUpload = true
				s.sendStartCode(codeStartUpload)
			}
			//进度
			s.SendToWs("3|" + strconv.Itoa(int(result.Upload.Progress*100)))
			full := strconv.Itoa(67 + int(result.Upload.Progress*33))
			s.SendToWs("4|" + full)
			//数据
			s.SendToWs("203|" + strconv.FormatFloat(result.Upload.Bandwidth, 'f', 0, 32))
			break
		case "result":
			if !s.IsFinish {
				s.IsFinish = true
				s.sendStartCode(codeFinish)
			}
			bytes, err := json.Marshal(result.Result)
			if err != nil {
				log.Println("格式化SpeedTest结果异常", err, result.Result)
			}
			s.SendToWs("204|" + string(bytes))
			//s.SendToWs()
			break
		}
	}

	s.IsStart = false
	s.IsStartPing = false
	s.IsStartDownload = false
	s.IsStartUpload = false
	s.IsFinish = false

	s.TimeOut = 60 * 5 //超过5分钟的话就强制终止
}

// 发送排队信息
func (s App_SpeedTest) sendListNum(num, all int, willSrart bool) {
	if willSrart {
		//开始测速
		//s.SendSourseTextToWs("5|" + s.TaskId + "|2|0") // + strconv.Itoa(num) + "|" + strconv.Itoa(all)
		s.SendToWs("2|0")
	} else {
		//排队中
		//s.SendSourseTextToWs("5|" + s.TaskId + "|2|1|" + strconv.Itoa(num) + "|" + strconv.Itoa(all))
		s.SendToWs("2|1|" + strconv.Itoa(num) + "|" + strconv.Itoa(all))
	}
}
func (s App_SpeedTest) sendStartCode(code int) {
	s.SendSourseTextToWs("5|" + s.TaskId + "|" + strconv.Itoa(code))
}

func (s App_SpeedTest) JoinTestList() {
	defer TestListLock.Unlock()
	TestListLock.Lock()

	TestList = append(TestList, &s)

	//加入队列
	s.SendSourseTextToWs("5|1|" + s.TaskId)
	//发送排队消息
	s.sendListNum(len(TestList), len(TestList), false)
}

func (s *App_SpeedTest) RunNow() {
	s.sendListNum(len(TestList), len(TestList), true)
	//if !s.IsStart {
	//	s.IsStart = true
	//	s.sendStartCode(codeStart)
	//}

	s.RunExternalPrograms()

	//发送结束命令
	s.SendToWs("0")

	defer func() {
		TestListLock.Unlock()
	}()
	TestListLock.Lock()

	if len(TestList) > 1 {
		TestList = TestList[1 : len(TestList)-2]
	} else {
		TestList = []*App_SpeedTest{}
	}

	//广播最新排队信息
	for i := 0; i < len(TestList); i++ {
		TestList[i].sendListNum(i+1, len(TestList), false)
	}
}
