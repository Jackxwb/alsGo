package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

type ApiReault struct {
	Code      int
	Success   bool
	Message   string
	Data      interface{}
	Timestamp time.Time
}

func (api ApiReault) String() string {
	marshal, err := json.Marshal(api)
	if err != nil {
		print(err)
	}
	return string(marshal)
}

func (api ApiReault) SetMessageNow(message string) ApiReault {
	api.Message = message
	return api
}
func (api ApiReault) SetDataNow(data interface{}) ApiReault {
	api.Data = data
	return api
}

func ApiReasultOk() ApiReault {
	return ApiReault{
		Code:      200,
		Success:   true,
		Timestamp: time.Now(),
	}
}
func ApiReasultFail(errMessage string) ApiReault {
	return ApiReault{
		Code:      501,
		Success:   false,
		Timestamp: time.Now(),
		Message:   errMessage,
	}
}

func (api ApiReault) ToContext(c *gin.Context) {
	c.JSON(api.Code, api)
}
func (api ApiReault) ToContextWS(c *websocket.Conn) {
	c.WriteJSON(api)
}

func (api ApiReault) OnlyData(c *gin.Context) {
	c.String(api.Code, api.Message, api.Data)
	fmt.Println("------")
	fmt.Println(api.Data)
	fmt.Println("------")
}
