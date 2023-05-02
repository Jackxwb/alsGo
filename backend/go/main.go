package main

import (
	"github.com/gin-gonic/gin"
	"golang/Config"
	"golang/controller"
	"golang/service/app"
)

func main() {
	r := InitSystem()

	//注册默认ws
	r.GET("/ws", func(c *gin.Context) {
		controller.DefWebSocket(c)
	})

	r.GET("speedtest/download", func(c *gin.Context) {
		app.JsStDownload(c)
	})
	r.POST("speedtest/upload", func(c *gin.Context) {
		app.JsStUpload(c)
	})
	//speedTest 静态文件
	r.GET("speedtest-static/:file", func(c *gin.Context) {

		fileName := c.Param("file")
		if fileName == "" {
			if len(Config.Config.BaseInfo.Testfiles) > 1 {
				fileName = Config.Config.BaseInfo.Testfiles[0] + ".test"
			}
		}
		c.File("public/speedtest-static/" + fileName)
	})

	ListenAddr := ":4000"
	if Config.Config.ListenAddr != "" {
		ListenAddr = Config.Config.ListenAddr
	} else {
		ListenAddr = ""
	}
	if Config.Config.Port != "" {
		ListenAddr += ":"
		ListenAddr += Config.Config.Port
	} else {
		ListenAddr += ":4000"
	}
	r.Run(ListenAddr)
}
