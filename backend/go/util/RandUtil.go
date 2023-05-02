package util

import (
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandomFloat() float64 {
	return r.Float64()
}
func GetRandomInt() int64 {
	return r.Int63n(10) - 5
}
func RandomInt(max int64) int64 {
	return r.Int63n(max)
}
