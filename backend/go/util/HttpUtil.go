package util

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type IpInfo struct {
	Ip      string
	Country string
	Local   string
}

// ReadZxincIpAddrIpv4 获取Ipv4地址
func ReadZxincIpAddrIpv4() (*IpInfo, error) {
	return readZxincIpAddr("http://v4.ip.zxinc.org/info.php?type=json")
}

// ReadZxincIpAddrIpv6 获取ipv6地址
func ReadZxincIpAddrIpv6() (*IpInfo, error) {
	return readZxincIpAddr("http://v6.ip.zxinc.org/info.php?type=json")
}

// readZxincIpAddr 在线获取地址
func readZxincIpAddr(url string) (*IpInfo, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered from panic<在线获取IP失败>:", err)
			// 处理 panic 错误
		}
	}()
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	apiResult := struct {
		Code int
		Data struct {
			Country string
			Local   string
			Myip    string
		}
	}{}
	err = json.Unmarshal(bytes, &apiResult)
	if err != nil {
		return nil, err
	}
	return &IpInfo{
		Ip:      apiResult.Data.Myip,
		Local:   apiResult.Data.Local,
		Country: apiResult.Data.Country,
	}, nil
}

func GetTrueIp(c *gin.Context) []string {
	ips := make([]string, 0)
	Try(func() {
		header := c.Request.Header
		ipHeader := strings.Split("X-Forwarded-For,Proxy-Client-IP,WL-Proxy-Client-IP,TTP_CLIENT_IP,HTTP_X_FORWARDED_FOR", ",")
		for _, ipH := range ipHeader {
			ip := header.Get(ipH)
			if ip != "" {
				//ips += ip + ","
				ips = append(ips, ip)
			}
		}
		ips = append(ips, c.Request.RemoteAddr)
		//ips += c.Request.RemoteAddr + "]"
	}, func(err interface{}) {})
	//ips += "]"
	return ips
}

func GetTrueIpString(c *gin.Context) string {
	ips := GetTrueIp(c)
	ip := ""
	for i := 0; i < len(ips); i++ {
		ip += ips[i]
		if i+1 < len(ips) {
			ip += ","
		}
	}
	return ip
}
