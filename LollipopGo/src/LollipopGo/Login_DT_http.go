package main

import (
	_ "LollipopGo/LollipopGo/player"
	"Proto"
	"Proto/Proto2"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/rpc/jsonrpc"
	"strconv"
)

/*
  登录服务器:
  http://127.0.0.1:8891/GolangLtdDT?Protocol=8&Protocol2=1
*/

func IndexHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		req.ParseForm()
		defer func() { // 必须要先声明defer，否则不能捕获到panic异常
			if err := recover(); err != nil {
				fmt.Println("%s", err)

				req.Body.Close()
			}
		}()
		Protocol, bProtocol := req.Form["Protocol"]
		Protocol2, bProtocol2 := req.Form["Protocol2"]

		if bProtocol && bProtocol2 {
			// 主协议判断
			if Protocol[0] == strconv.Itoa(Proto.G_GameLogin_Proto) {
				// 子协议判断
				switch Protocol2[0] {
				case strconv.Itoa(Proto2.C2GL_GameLoginProto2):
					// DB server 获取 验证信息  rpc 操作
					//------------------------------------------------------
					// 暂时不解析用户名和密码 --> 后面独立出来再增加！！！
					data := DB_rpc_()
					b, _ := json.Marshal(data)
					fmt.Fprint(w, base64.StdEncoding.EncodeToString(b))
					//------------------------------------------------------
					return
				default:
					fmt.Fprintln(w, "88902")
					return
				}
			}
			fmt.Fprintln(w, "88904")
			return
		}
		// 服务器获取通信方式错误 --> 8890 + 1
		fmt.Fprintln(w, "88901")
		return
	}

}

// jsonrpc 数据处理
func DB_rpc_() interface{} {
	// 链接DB操作
	client, err := jsonrpc.Dial("tcp", service)
	if err != nil {
		fmt.Println("dial error:", err)
	}
	// 测试 --
	args := Args{1, 2}
	// 返回数据的结构体 -->  消息的结构
	// 正常是读取数据库后得到的
	var reply Proto2.GL2C_GameLogin
	// 同步调用
	err = client.Call("Arith.Muliply", args, &reply)
	if err != nil {
		fmt.Println("Arith.Muliply call error:", err)
	}

	// divCall := client.Go("Arith.Muliply", args, &reply, nil)
	// replyCall := <-divCall.Done // will be equal to divCall
	// fmt.Println(replyCall.Reply)
	// 返回的数据
	fmt.Println("the arith.mutiply is :", reply)
	return reply
}
