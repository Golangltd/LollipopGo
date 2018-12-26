package main

import (
	_ "LollipopGo/LollipopGo/player"
	"Proto"
	"Proto/Proto2"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/rpc/jsonrpc"
	"strconv"
)

/*
  Gm 游戏服务器：
	1 修改游戏服务器中的玩家的个人的数据的变化，例如：金币，M卡等
	2 玩家等级的限制
	3 协议处理
*/

var ConnRPC_GM *rpc.Client // 保存全局数据

// 初始化RPC
func init {
	client, err := jsonrpc.Dial("tcp", service)
	if err != nil {
		log.Debug("dial error:", err)
		return
	}
	ConnRPC_GM = client
}

func IndexHandlerGM(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		req.ParseForm()
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("%s", err)

				req.Body.Close()
			}
		}()

		Protocol, bProtocol := req.Form["Protocol"]
		Protocol2, bProtocol2 := req.Form["Protocol2"]
		if bProtocol && bProtocol2 {
			if Protocol[0] == strconv.Itoa(Proto.G_GameGM_Proto) {
				switch Protocol2[0] {
				case strconv.Itoa(Proto2.W2GMS_Modify_PlayerDataProto2):
					// 修改数据玩家的结构的数据信息
					// DB server 获取 rpc 操作
					//------------------------------------------------------
					// 修改Gm数据
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


// 数据操作
func ModefyGamePlayerData(){
	
	return
}