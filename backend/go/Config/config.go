package Config

import (
	"fmt"
	"github.com/jinzhu/configor"
)

var Config = struct {
	ListenAddr string  `default:"0.0.0.0" json:"listenaddr"`
	Port       string  `default:"4000" json:"port"`
	Testfiles  []int64 `default:"[1048576,10485760,104857600]" json:"testfiles"` //,1073741824 (1G)
	BaseInfo   struct {
		PublicIpv4               string   `json:"public_ipv4"`
		PublicIpv6               string   `json:"public_ipv6"`
		Location                 string   `json:"location"`
		Bandwidth                bool     ` json:"bandwidth"  default:"false"`
		DisplayTraffic           bool     ` json:"display_traffic" default:"true"`
		DisplaySpeedtest         bool     ` json:"display_speedtest" default:"true"`
		UtilitiesPing            bool     ` json:"utilities_ping" default:"true"`
		UtilitiesTraceroute      bool     ` json:"utilities_traceroute" default:"true"`
		UtilitiesIperf3          bool     ` json:"utilities_iperf3" default:"true"`
		UtilitiesIperf3PortMin   string   ` json:"utilities_iperf3_port_min" default:"30000"`
		UtilitiesIperf3PortMax   string   ` json:"utilities_iperf3_port_max" default:"31000"`
		UtilitiesSpeedtestdotnet bool     ` json:"utilities_speedtestdotnet" default:"true"`
		UtilitiesFakeshell       bool     ` json:"utilities_fakeshell" default:"false"`
		SponsorMessage           string   `json:"sponsor_message"`
		Testfiles                []string `json:"testfiles"`
		ClientIp                 string   `json:"client_ip"`
	} `json:"baseinfo"`
	App struct {
		DefTimeOut         int    `json:"deftimeout" default:"90"`
		PingResultTemplate string `default:"" json:"pingresulttemplate"`
		// ping 超时
		PingResultNotLinkTemplate string `default:"" json:"pingresultnotlinktemplate"`
		// ping 找不到目标
		PingResultNotPingTemplate string `default:"" json:"pingresultnotpingtemplate"`
		//固定 Iperf 端口
		IperfFixedPort int `json:"iperffixedport" default:0`
	} `json:"app"`
	NetSet struct {
		NetworkAdapter struct {
			//只显示部分 网络适配器
			OnlyShow []string `json:"onlyshow"`
			//不显示的 网络适配器
			DontShow []string `json:"dontshow"`
		} `json:"networkadapter"`
		//是否隐藏外网IP地址
		HideExternalIP bool `json:"hideexternalip" default:"false"`
	} `json:"netset"`
}{}

func LoadYMLConfig() {
	err := configor.Load(&Config, "config.yml")
	if err != nil {
		fmt.Println("load yml config Err:", err)
	}
}
