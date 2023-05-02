package app

import (
	"github.com/gorilla/websocket"
	"golang/Config"
	"golang/util"
	"log"
	"runtime"
	"strconv"
)

type App_Ping struct {
	BaseRun
	os                      string
	pingNum                 int
	resultTemplate          string
	resultTemplateTranslate string
	//无法连接到目标的模板
	resultNotLinkTemplate          string
	resultNotLinkTemplateTranslate string
	//无法解析目标
	resultNotPingTemplate          string
	resultNotPingTemplateTranslate string
	//内置，找不到就用内置的，来解决window的兼容问题
	seq           int
	useBuiltInSeq bool
	//传入目标
	inHost string
}

func (p App_Ping) sendFailMessage() {
	p.SendToWs("-1")
}
func (p App_Ping) sendErrMessage(txt string) {
	p.SendToWs(txt + "|0|-|-")
}

func (p *App_Ping) InitPing(conn *websocket.Conn, tager string) {
	p.Init()

	p.Exe = "ping"
	p.pingNum = 10
	p.Ws = conn

	p.SendTTextPrefix = "1|" + p.TaskId + "|1|"
	p.seq = 1
	p.useBuiltInSeq = false
	p.inHost = tager

	p.os = runtime.GOOS
	switch p.os {
	case "windows":
		p.Arg = []string{"-n", strconv.Itoa(p.pingNum), p.inHost}
		p.resultTemplate = "来自 (?P<ip>.*?) 的回复: 字节=.*? 时间[=<]?(?P<time>.*?)ms TTL=(?P<ttl>.*)$"
		p.resultNotLinkTemplate = "请求超时。"
		p.resultNotPingTemplate = "找不到主机"
		break
	default:
		p.Arg = []string{"-O", "-c", strconv.Itoa(p.pingNum), p.inHost}
		p.resultTemplate = "from (?P<ip>.*?) \\(.*\\): icmp_seq=(?P<seq>.*?) ttl=(?P<ttl>.*) time=(?P<time>.*?) ms"
		p.resultNotLinkTemplate = "no answer yet for icmp_seq="
		//p.resultNotPingTemplate = "未知的名称或服务"
		p.resultNotPingTemplate = "bad address"
		break
	}
	//用户可以外部自行配置模板
	p.resultTemplateTranslate = Config.Config.App.PingResultTemplate
	p.resultNotLinkTemplateTranslate = Config.Config.App.PingResultNotLinkTemplate
	p.resultNotPingTemplateTranslate = Config.Config.App.PingResultNotPingTemplate

	p.DoText = func(txt string) {
		//log.Println(">>", txt)
		//正常数据
		find := util.RegexpFindStringGroup(p.resultTemplateTranslate, txt)
		if len(find) == 0 { // find == nil ||
			find = util.RegexpFindStringGroup(p.resultTemplate, txt)
		}
		if len(find) > 0 { // find != nil ||
			if !p.useBuiltInSeq {
				//使用回传seq
				if find["seq"] != "" {
					atoi, err := strconv.Atoi(find["seq"])
					if err != nil {
						log.Println(err)
						p.useBuiltInSeq = true
					} else {
						//外部seq赋值
						p.seq = atoi
					}
				} else {
					p.useBuiltInSeq = true
				}
			}
			p.SendToWs(find["ip"] + "|" + strconv.Itoa(p.seq) + "|" + find["ttl"] + "|" + find["time"])
			if p.useBuiltInSeq {
				p.seq++
			}
			return
		}
		//无法ping到目标模板
		cannotPing := false
		if p.resultNotLinkTemplateTranslate != "" {
			cannotPing = util.RegexpFindString(p.resultNotLinkTemplateTranslate, txt)
		}
		if !cannotPing {
			cannotPing = util.RegexpFindString(p.resultNotLinkTemplate, txt)
		}
		if cannotPing {
			p.sendFailMessage()
			return
		}

		//无法解析目标
		badAddress := false
		if p.resultNotPingTemplateTranslate != "" {
			badAddress = util.RegexpFindString(p.resultNotPingTemplateTranslate, txt)
		}
		if !badAddress {
			badAddress = util.RegexpFindString(p.resultNotPingTemplate, txt)
		}
		if badAddress {
			//p.sendErrMessage("解析不到主机")
			p.sendErrMessage(txt)
			return
		}

		//都匹配不到就不处理
		//log.Println("xx", find, cannotPing, txt)
	}
}

func (p App_Ping) RunPing() {
	p.SendSourseTextToWs("1|1|" + p.inHost + "|" + p.TaskId)

	p.RunExternalPrograms()

	p.SendSourseTextToWs("1|" + p.TaskId + "|0")
}
