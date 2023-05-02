package app

import (
	"github.com/gin-gonic/gin"
	"golang/util"
	"io"
	"runtime"
	"strconv"
)

func JsStDownload(c *gin.Context) {
	defer func() {
		runtime.GC()
	}()
	query := c.Query("ckSize")
	ckSize, err := strconv.ParseInt(query, 10, 32)
	if err != nil {
		ckSize = 100
	}
	if ckSize > 50 {
		ckSize = 50
	}
	ckSize += util.GetRandomInt()
	//c.Writer.Write(RequestMemory(int64(float64(ckSize) * 1024 * 1024 * GetRandomFloat())))
	c.Writer.Write(RequestMemory(int64(float64(ckSize << 20))))
}
func RequestMemory(size int64) []byte {
	bytes := make([]byte, size)
	for i := range bytes {
		if i == 0 {
			bytes[i] = 0
		} else {
			bytes[i] = bytes[i-1] + 1
		}
		//bytes[i] = byte(RandomInt(255))
	}
	return bytes
}

func JsStUpload(c *gin.Context) {
	_, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(501, err.Error())
	}
	c.String(200, "")
}
