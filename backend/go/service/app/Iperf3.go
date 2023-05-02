package app

import (
	"github.com/gorilla/websocket"
	"golang/Config"
	"golang/util"
	"strconv"
)

type App_Iperf3 struct {
	BaseRun
	openProt int
}

func (i *App_Iperf3) InitIperf3(Ws *websocket.Conn) error {
	//log.Println("父对象初始化")
	i.Init()
	i.Exe = "./app/iperf3"

	//log.Println("赋值ws")
	i.Ws = Ws
	//log.Println("ws -> ", i.Ws)

	if !i.findFileIsExist(i.Exe) {
		if i.findFileIsExist("./app/iperf3.exe") {
			i.Exe = "./app/iperf3.exe"
		} else {
			//errText := "iperf3 程序丢失，请确保其在 ./app/iperf3"
			////i.SendToWs(errText)
			//return errors.New(errText)
			i.Exe = "iperf3"
		}
	}

	i.SendTTextPrefix = "4|" + i.TaskId + "|2|"

	if Config.Config.App.IperfFixedPort > 0 {
		i.openProt = Config.Config.App.IperfFixedPort
	} else {
		min, err := strconv.Atoi(Config.Config.BaseInfo.UtilitiesIperf3PortMin)
		if err != nil {
			min = 30000
		}
		max, err := strconv.Atoi(Config.Config.BaseInfo.UtilitiesIperf3PortMax)
		if err != nil {
			max = 31000
		}
		i.openProt = int(int64(min) + util.RandomInt(int64(max-min)))
	}

	i.Arg = []string{"-s", "-p", strconv.Itoa(i.openProt)}

	i.DoText = func(txt string) {
		i.SendToWs(txt + "\r\n")
	}

	return nil
}

func (i App_Iperf3) Run() {
	i.SendSourseTextToWs("4|1|" + i.TaskId)
	//str := strconv.FormatFloat(i.TimeOut.Seconds(), 'f', 0, 32)
	//i.SendSourseTextToWs("4|" + i.TaskId + "|1|" + strconv.Itoa(i.openProt) + "|" + str)
	i.SendSourseTextToWs("4|" + i.TaskId + "|1|" + strconv.Itoa(i.openProt) + "|" + strconv.Itoa(i.TimeOut))

	i.RunExternalPrograms()
	//结束信息
	i.SendSourseTextToWs("4|" + i.TaskId + "|0")
}
