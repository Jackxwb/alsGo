package service

import "net/http"
import _ "net/http/pprof"

func CreatePpoofServie() {
	go func() {
		http.ListenAndServe("0.0.0.0:8889", nil)
	}()
}
