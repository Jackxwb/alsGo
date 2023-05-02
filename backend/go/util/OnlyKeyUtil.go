package util

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

var gen *snowflake.Node
var isInit = false

func init_() {
	node, err := snowflake.NewNode(4000 % 1023)
	if err != nil {
		log.Println("snowflake 初始化失败!", err)
	}
	gen = node
	isInit = true
}

func OnlyKey() string {
	if !isInit {
		init_()
	}
	return gen.Generate().String()
}
func OnlyKeyInt() int64 {
	if !isInit {
		init_()
	}
	return gen.Generate().Int64()
}
