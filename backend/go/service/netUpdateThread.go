package service

import (
	"fmt"
	goNet "github.com/shirou/gopsutil/net"
	"golang/Config"
	"golang/util"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// 缓存网络帧
var netDataFrame []*DataFrame

// 锁
var lock sync.RWMutex

// 返回外部
func GetNetDataFrame() []*DataFrame {
	defer func() {
		lock.RUnlock()
	}()
	lock.RLock()

	return netDataFrame
}

type DataFrame struct {
	Time int64          `json:"time"`
	Data []*NetWorkInfo `json:"data"`
}

// 回传实时数据
type NetWorkInfo struct {
	Name string `json:"name"`
	Recv string `json:"recv"`
	Send string `json:"send"`
}

//func (df *DataFrame) cleanData() {
//	df.Time = 0
//	for i := 0; i < len(df.Data); i++ {
//		df.Data[i] = nil
//	}
//	df.Data = nil
//}

// 内部网卡当前数据
var netWorkData = make(map[string]*NetWorkInfo)

// 返回外部
func GetNetWorkNowData() map[string]*NetWorkInfo {
	defer func() {
		lock.RUnlock()
	}()
	lock.RLock()

	return netWorkData
}

var NetUpdateThreadIsRun = false

func NetUpdateThread() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered from panic:", err)
			// 处理 panic 错误
		}
		log.Println("网络更新服务异常退出！！")
		NetUpdateThreadIsRun = false
	}()

	//debug := true
	//if debug {
	//	return
	//}

	NetUpdateThreadIsRun = true
	for true {
		//fmt.Println("---update---")
		//newWork()
		newWorkV2()
		//fmt.Println("---update ed---")
		//time.Sleep(time.Second * 1)

		// Debug ram
		debugRam()
		// Debug ram

		time.Sleep(time.Millisecond * 500)
	}
}

var ramMax, ramMin uint64

var lastDebugTime int64

func debugRam() {
	//避免高频写出硬盘
	now := time.Now().Unix()
	if now-lastDebugTime < 10 {
		return
	}
	lastDebugTime = now

	ram := runtime.MemStats{}
	runtime.ReadMemStats(&ram)
	allocNow := ram.Alloc
	if allocNow > ramMax {
		ramMax = allocNow
	}
	if ramMin > allocNow || ramMin == 0 {
		ramMin = allocNow
	}
	//log.Println("内存Alloc:", util.FomatSize(float64(allocNow)), ", 最大:", util.FomatSize(float64(ramMax)), "最小:", util.FomatSize(float64(ramMin)), "在线:", GetOnlineNum())
	//data := "内存Alloc:" + util.FomatSize(float64(allocNow)) + ", 最大:" + util.FomatSize(float64(ramMax)) + "最小:" + util.FomatSize(float64(ramMin)) + "在线:" + GetOnlineNum()
	data := fmt.Sprintf("%s 内存Alloc: %s , 最大: %s 最小: %s 在线: %d\r\n",
		util.GetNowStr(),
		util.FomatSize(float64(allocNow)),
		util.FomatSize(float64(ramMax)),
		util.FomatSize(float64(ramMin)),
		GetOnlineNum(),
	)

	onlineUser := GetOnlineUser()
	for i := 0; i < len(OnlineList); i++ {
		data += fmt.Sprintf("->\t%d [%s](%s) [E:%d]%s\r\n",
			i,
			onlineUser[i].Addr,
			onlineUser[i].Ips,
			onlineUser[i].ErrNum,
			onlineUser[i].LastErr,
		)
	}

	err := util.SaveText([]byte(data), "./onlineInfo")
	if err != nil {
		log.Println("Debug异常", err)
	}
}

// 统计流量
func newWork() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[newWork]Recovered from panic:", err)
			// 处理 panic 错误
		}
		lock.Unlock()
	}()
	lock.Lock()

	frame := DataFrame{Time: time.Now().Unix()}

	networkIOCount, _ := goNet.IOCounters(true)
	for _, v := range networkIOCount {
		//跳过无返回的网卡
		if v.BytesRecv == 0 {
			continue
		}

		//用户自己的配置
		if !canShow(v.Name) {
			continue
		}

		workInfo := netWorkData[v.Name]
		if workInfo == nil {
			workInfo = &NetWorkInfo{
				Name: v.Name,
			}
		}

		workInfo.Recv = strconv.FormatUint(v.BytesRecv, 10)
		workInfo.Send = strconv.FormatUint(v.BytesSent, 10)

		netWorkData[v.Name] = workInfo

		frame.Data = append(frame.Data, workInfo)
	}

	if len(netDataFrame) > 10 {
		//netDataFrame[0].cleanData()
		netDataFrame[0] = nil
		netDataFrame = netDataFrame[1 : len(netDataFrame)-1]
	}
	netDataFrame = append(netDataFrame, &frame)

	//触发发送
	UpdateNetMessage(frame)
}

// 统计流量
func newWorkV2() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("[newWorkV2]Recovered from panic:", err)
			// 处理 panic 错误
		}
		lock.Unlock()
	}()
	lock.Lock()

	last := len(netDataFrame)
	if last > 10 {
		for i := 1; i < last-1; i++ {
			netDataFrame[i-1] = netDataFrame[i]
		}
	} else {
		netDataFrame = append(netDataFrame, &DataFrame{Time: time.Now().Unix()})
		last++
	}

	networkIOCount, _ := goNet.IOCounters(true)
	for _, v := range networkIOCount {
		//跳过无返回的网卡
		if v.BytesRecv == 0 {
			continue
		}

		//用户自己的配置
		if !canShow(v.Name) {
			continue
		}

		workInfo := netWorkData[v.Name]
		if workInfo == nil {
			workInfo = &NetWorkInfo{
				Name: v.Name,
			}
		}

		workInfo.Recv = strconv.FormatUint(v.BytesRecv, 10)
		workInfo.Send = strconv.FormatUint(v.BytesSent, 10)

		netWorkData[v.Name] = workInfo

		//log.Println("debug [network]", workInfo.Name, workInfo.Recv, workInfo.Send)

		//frame.Data = append(frame.Data, workInfo)
		netDataFrame[last-1].Data = append(netDataFrame[last-1].Data, workInfo)
	}

	frame := netDataFrame[last-1]

	//触发发送
	UpdateNetMessage(*frame)
}

func canShow(name string) bool {
	for _, item := range Config.Config.NetSet.NetworkAdapter.DontShow {
		if util.RegexpFindString(item, name) {
			return false
		}
	}
	for _, item := range Config.Config.NetSet.NetworkAdapter.OnlyShow {
		if util.RegexpFindString(item, name) {
			return true
		}
	}
	//是否开启只显示白名单
	if len(Config.Config.NetSet.NetworkAdapter.OnlyShow) > 0 {
		return false
	}
	return true
}

func GetNetworkAdaptersWithTraffic() map[string]bool {
	netAdapters := make(map[string]bool)
	networkIOCount, _ := goNet.IOCounters(true)
	for _, v := range networkIOCount {
		//跳过无返回的网卡
		if v.BytesRecv == 0 {
			continue
		}

		//用户自己的配置
		if !canShow(v.Name) {
			continue
		}

		//netAdapters = append(netAdapters, v.Name)
		netAdapters[v.Name] = true
	}
	return netAdapters
}
