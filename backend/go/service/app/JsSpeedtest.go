package app

import (
	"github.com/gin-gonic/gin"
	"golang/util"
	"io"
	"runtime"
	"strconv"
	"sync"
)

var memoryBuffData []byte
var memoryBuffDataLock sync.RWMutex

// CleanMemoryBuffData 清理测速缓存
func CleanMemoryBuffData() {
	defer memoryBuffDataLock.Unlock()
	memoryBuffDataLock.Lock()

	for i := 0; i < len(memoryBuffData); i++ {
		memoryBuffData[i] = 0
	}
	memoryBuffData = nil
}

// 创建测速缓存
func checkMemoryBuffData(ckSize int64) {
	defer memoryBuffDataLock.Unlock()
	memoryBuffDataLock.Lock()

	if memoryBuffData == nil {
		memoryBuffData = RequestMemory(int64(float64(ckSize << 20)))
	}
}

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
	//c.Writer.Write(RequestMemory(int64(float64(ckSize << 20))))

	//缓存不存在就创建，存在就跳过使用缓存
	checkMemoryBuffData(ckSize)

	c.Writer.Write(memoryBuffData)
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
