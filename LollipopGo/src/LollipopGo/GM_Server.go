package main

import (
	"LollipopGo/LollipopGo/log"
	"LollipopGo/LollipopGo/player"
	"Proto"
	"Proto/Proto2"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
)

/*
  Gm 游戏服务器：
	1 修改游戏服务器中的玩家的个人的数据的变化，例如：金币，M卡等
	2 玩家等级的限制
	3 协议处理
*/

var ConnRPC_GM *rpc.Client // 保存全局数据
// 链接 http://127.0.0.1:8892/GolangLtdGM?Protocol=11&Protocol2=1&UID=&Itype=&ModifyNum=

// 初始化RPC
func init() {
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
					UID, bUID := req.Form["UID"]
					Itype, bItype := req.Form["Itype"]
					ModifyNum, bModifyNum := req.Form["ModifyNum"]
					// 修改Gm数据
					if bUID && bItype && bModifyNum {
						data := ModefyGamePlayerData(UID[0], Itype[0], ModifyNum[0])
						b, _ := json.Marshal(data)
						fmt.Fprint(w, base64.StdEncoding.EncodeToString(b))
						//------------------------------------------------------
						return
					} else {
						fmt.Fprint(w, base64.StdEncoding.EncodeToString([]byte("参数错误!")))
						//------------------------------------------------------
						return
					}
				case strconv.Itoa(Proto2.W2GMS_Modify_PlayerEmailDataProto2):
					{ // 跑马灯+邮件的协议
						IMsgtype, bIMsgtype := req.Form["IMsgtype"]
						if bIMsgtype {

							if IMsgtype[0] == "1" {
								// 邮件相关
								EmailData, bEmailData := req.Form["EmailData"]
								if bEmailData {
									fmt.Println("EmailData", EmailData[0])
									stb := &player.EmailST{}
									json.Unmarshal([]byte(EmailData[0]), &stb)
									ModefyGameEmailData(stb)
									fmt.Fprint(w, base64.StdEncoding.EncodeToString([]byte("true")))
									return
								}
							} else if IMsgtype[0] == "2" {
								// 跑马灯
								// MsgData, bMsgData := req.Form["MsgData"]
								// if bMsgData {
								// }
							}
						}
					}
				default:
					fmt.Println(Protocol2[0])
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

//------------------------------------------------------------------------------
// 邮件通知
func ModefyGameEmailData(data *player.EmailST) interface{} {

	// 返回的数据
	var reply Proto2.GMS2W_Modify_PlayerEmailData

	fmt.Println("data:", data)
	// 异步调用
	divCall := ConnRPC_GM.Go("Arith.ModefyPlayerEmailDataGM", data, &reply, nil)
	replyCall := <-divCall.Done
	fmt.Println(replyCall.Reply)
	fmt.Println("the ModefyGameEmailData is :", reply)
	return reply
}

//------------------------------------------------------------------------------
// GM 修改数据
func ModefyGamePlayerData(uid, itype, modifynum string) interface{} {
	// 发送的数据
	args := Proto2.W2GMS_Modify_PlayerData{
		UID:       uid,
		Itype:     itype,
		ModifyNum: modifynum,
	}
	// 返回的数据
	var reply Proto2.GMS2W_Modify_PlayerData
	//--------------------------------------------------------------------------
	// 同步调用
	// err = ConnRPC_GM.Call("Arith.ModefyPlayerDataGM", args, &reply)
	// if err != nil {
	// 	fmt.Println("Arith.ModefyPlayerDataGM call error:", err)
	// }
	// 异步调用
	divCall := ConnRPC_GM.Go("Arith.ModefyPlayerDataGM", args, &reply, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	fmt.Println(replyCall.Reply)
	//--------------------------------------------------------------------------
	// 返回的数据
	fmt.Println("the arith.mutiply is :", reply)
	return reply
}
