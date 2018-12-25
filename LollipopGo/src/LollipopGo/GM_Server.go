package main

import (
	"fmt"
	"net/http"
)

/*
  Gm 游戏服务器：
	1 修改游戏服务器中的玩家的个人的数据的变化，例如：金币，M卡等
	2 玩家等级的限制
*/

func IndexHandlerGM(w http.ResponseWriter, r *http.Request) {
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

	fmt.Fprintln(w, "请用Get 方式请求!")
	return
}
