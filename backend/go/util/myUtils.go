package util

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Try try catch 简单实现
func Try(function func(), catch func(err interface{})) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Errorf("TtyErr: %v", err)
			catch(err)
		}
	}()

	function()
}

var (
	dw      = strings.Split("B,KB,MB,GB,TB,PB,EB,ZB,YB", ",")
	logBuff = math.Log(1024)
)

//func FomatSize(size float64) string {
//	tmp := size
//	for i := 0; i < len(dw); i++ {
//		if tmp < 1024 {
//			return fmt.Sprintf("%.2f %s", tmp, dw[i])
//		}
//		tmp = tmp / 1024
//	}
//	return fmt.Sprintf("%.2f %s", tmp, dw[len(dw)-1])
//}

// FomatSizeP FomatSize 第二版实现
func FomatSizeP(size float64, prec int) string {
	ii := int(math.Floor(math.Log(size) / logBuff))
	if ii >= len(dw) {
		ii = len(dw) - 1
	}
	buff := 1 << uint64(ii*10)
	//return fmt.Sprintf("%.2f %s", size/float64(buff), dw[ii])
	val := math.Round(size*100/float64(buff)) / 100
	valStr := strconv.FormatFloat(val, 'f', prec, 64)
	return valStr + dw[ii]
}
func FomatSize(size float64) string {
	return FomatSizeP(size, 2)
}
