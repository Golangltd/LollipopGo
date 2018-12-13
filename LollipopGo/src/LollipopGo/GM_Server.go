package main

import (
	"fmt"
	"net/http"
)

/*
  登录服务器的函数
*/

func IndexHandlerGM(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
	// 需要处理 get请求等
}
