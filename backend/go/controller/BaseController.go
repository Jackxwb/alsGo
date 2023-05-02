package controller

import (
	"github.com/gin-gonic/gin"
)

func NoFindResult(c *gin.Context) {
	//c.HTML(404, "404.html", gin.H{})
	//url := c.Request.URL
	//fmt.Println(url)
	c.String(404, "404 no find")
}
func NoAuthority(c *gin.Context) {
	c.String(403, "403 noAuthority")
}
