package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nanmu42/gzip"
	"golang/Config"
	"golang/ScheduledTask"
	"golang/controller"
	"golang/service"
	"golang/service/app"
	"golang/util"
	"log"
	"net"
	"os"
	"os/user"
)

func InitSystem() *gin.Engine {

	//内存分析
	//service.CreatePpoofServie()

	//打印基础信息
	//printBaseInfo()

	//读取配置文件
	Config.LoadYMLConfig()

	//获取ip地址
	initServerIp()

	//定时任务
	ScheduledTask.InitCornTask()

	//自动更新网卡流量
	go service.NetUpdateThread()

	//初始化静态测速文件
	initSpeedtestStaticFile()

	return initGin()
}

func printBaseInfo() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println("获取启动路径失败", err)
	}
	log.Println("程序工作在->", dir)

	u, err := user.Current()
	if err != nil {
		log.Println("获取运行用户失败", err)
	}
	log.Println("当前用户->", u.Uid, u.Gid, u.Name, u.Username, u.HomeDir)

	//删除环境变量
	//delEnv()

	environ := os.Environ()
	//log.Println("Env->", environ)
	log.Println("Env->")
	for name, val := range environ {
		fmt.Println(name, "=", val)
	}
	log.Println("<--Env")
}

//func delEnv() {
//	removeEnvTest := []string{"LC_ALL", "LS_COLORS", "LC_MEASUREMENT", "SSH_CONNECTION", "LESSCLOSE", "LC_PAPER", "LC_MONETARY", "LANG", "DISPLAY", "LC_NAME", "SSH_TTY", "MAIL", "TERM", "SHELL", "SHLVL", "LANGUAGE", "LC_TELEPHONE", "LOGNAME", "XDG_RUNTIME_DIR", "PATH", "LC_IDENTIFICATION", "LESSOPEN", "LE_WORKING_DIR", "LC_TIME", "_", "OLDPWD"}
//	for _, key := range removeEnvTest {
//		unRegEnv(key)
//	}
//}

//func unRegEnv(name string) {
//	os.Unsetenv(name)
//}

func initGin() *gin.Engine {
	r := gin.New()
	//403
	r.GET("/system/403", controller.NoAuthority)
	//404
	r.NoRoute(controller.NoFindResult)

	//前端文件
	r.Static("assets", "public/static/dist/assets")
	r.GET("/", func(c *gin.Context) {
		c.File("public/static/dist/index.html")
	})
	r.GET("favicon.ico", func(c *gin.Context) {
		c.File("public/static/dist/favicon.ico")
	})

	r.GET("speedtest_worker.js", func(c *gin.Context) {
		c.File("public/static/speedtest_worker.js")
	})

	//gzip压缩
	r.Use(gzip.DefaultHandler().Gin)

	return r
}

func initServerIp() {
	//是否隐藏外网IP地址
	if Config.Config.NetSet.HideExternalIP {
		getLocalIp()
		return
	}

	ipv4, err := util.ReadZxincIpAddrIpv4()
	if err != nil {
		log.Println("获取IPv4错误 -", err)
		//获取失败的情况下获取本地地址
		getLocalIp()
	} else {
		Config.Config.BaseInfo.PublicIpv4 = ipv4.Ip
		Config.Config.BaseInfo.Location = ipv4.Local + " " + ipv4.Country
	}
	ipv6, err := util.ReadZxincIpAddrIpv6()
	if err != nil {
		log.Println("获取IPv6错误 -", err)
	} else {
		//webSocketService.BaseInfo.PublicIpv6 = ipv6.Ip
		Config.Config.BaseInfo.PublicIpv6 = ipv6.Ip
	}
}

func getLocalIp() {
	Config.Config.BaseInfo.Location = "局域网"
	Config.Config.BaseInfo.PublicIpv4 = ""

	networkAdaptersNames := service.GetNetworkAdaptersWithTraffic()

	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {
		if networkAdaptersNames[interf.Name] {
			addrs, err := interf.Addrs()
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, addr := range addrs {
				// 检查ip地址判断是否回环地址
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsLinkLocalUnicast() {
					if ipnet.IP.To4().String() == "<nil>" {
						//ipv6 考虑到机器可能没有外网ipv6的情况，暂不处理
					} else {
						//ipv4
						Config.Config.BaseInfo.PublicIpv4 = Config.Config.BaseInfo.PublicIpv4 + ipnet.IP.String() + "; "
					}
				}
			}
		}
	}
}

func initSpeedtestStaticFile() {
	//路径
	savePath := "public/speedtest-static"

	//目录是否存在
	_, err := os.Open(savePath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(savePath, 0766)
		} else {
			log.Println("创建静态测速文件失败", err.Error())
			return
		}
	}
	//fmt.Println(Config.Config.Testfiles)
	log.Println("初始化静态测速文件: ", Config.Config.Testfiles)

	for _, size := range Config.Config.Testfiles {
		fileName := util.FomatSizeP(float64(size), 0)
		log.Print("->", fileName)
		fullFile := savePath + "/" + fileName + ".test"
		open, err := os.Open(fullFile)
		if err != nil {
			if os.IsNotExist(err) {
				os.Create(fullFile)
			}
		} else {
			//os.Remove(fullFile)
			Config.Config.BaseInfo.Testfiles = append(Config.Config.BaseInfo.Testfiles, fileName)
			log.Println(" | 文件已存在，跳过")
			continue
		}
		log.Println("")
		defer func() {
			open.Close()
		}()

		file, err := os.OpenFile(fullFile, os.O_WRONLY|os.O_CREATE, 0766)
		if err != nil {
			fmt.Println(err)
			return
		}

		blockSize := int64(1024 * 1024)

		lost := size
		for lost > blockSize {
			randomInt := util.RandomInt(blockSize)
			file.Write(app.RequestMemory(randomInt))
			lost -= randomInt
		}
		if lost > 0 {
			file.Write(app.RequestMemory(lost))
		}
		Config.Config.BaseInfo.Testfiles = append(Config.Config.BaseInfo.Testfiles, fileName)
	}
	log.Println("静态测速文件初始化完成", Config.Config.BaseInfo.Testfiles)
}
