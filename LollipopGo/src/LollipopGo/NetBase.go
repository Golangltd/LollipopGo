package main

import (
	"encoding/json"
	"fmt"

	"Proto"
	"Proto/Proto2"

	"code.google.com/p/go.net/websocket"
)

func wwwGolangLtd(ws *websocket.Conn) {
	// fmt.Println("Golang语言社区 欢迎您！", ws)
	// data = json{}
	data := ws.Request().URL.Query().Get("data")
	fmt.Println("data:", data)
	NetDataConntmp := &NetDataConn{
		Connection:    ws,
		StrMd5:        "",
		MapSafe:       M,
		MapSafeServer: MServer,
	}
	NetDataConntmp.PullFromClient()
}

// 公用的send函数
func PlayerSendToServer(conn *websocket.Conn, data interface{}) {
	jsons, err := json.Marshal(data)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	///fmt.Println("jsons:", string(jsons))
	errq := websocket.Message.Send(conn, jsons)
	if errq != nil {
		fmt.Println(errq)
	}
	return
}

// 发送给客户端的数据信息函数
func (this *NetDataConn) PlayerSendMessage(senddata interface{}) {
	// 1 消息序列化 interface --->  json
	b, errjson := json.Marshal(senddata)
	if errjson != nil {
		fmt.Println(errjson.Error())
		return
	}
	// 数据转换 json的byte数组 --->  string
	data := "data:" + string(b[0:len(b)])
	fmt.Println(data)
	// 发送
	err := websocket.JSON.Send(this.Connection, b)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

// 群发广播函数:同一个房间
func PlayerSendBroadcastToRoomPlayer(iroomID int) {

	//  处理数据操纵
	// -------------------------------------------------------------------------
	for itr := M.Iterator(); itr.HasNext(); {
		k, v, _ := itr.Next()
		strsplit := Strings_Split(k.(string), "|")
		for i := 0; i < len(strsplit); i++ {
			if len(strsplit) < 2 {
				continue
			}

			if strsplit[2] == "room" {
				// 进行数据的查询类型
				switch v.(interface{}).(type) {
				case *NetDataConn:
					{ // 发送数据操作
						data := &Proto2.C2S_PlayerAddGame{
							Protocol:      Proto.GameNet_Proto,
							Protocol2:     Proto2.Net_Kicking_PlayerProto2,
							OpenID:        "1212334",
							RoomID:        1,
							PlayerHeadURL: "11",
							Init_X:        13,
							Init_Y:        10,
						}
						// 发送数据
						v.(interface{}).(*NetDataConn).PlayerSendMessage(data)
					}
				}
			}
		}
	}
	// -------------------------------------------------------------------------
}

// 推送给服务器
// 推送给 Global server key = Global_Server
// 几个小游戏的serverID
func (this *NetDataConn) SendServerDataFunc(StrMD5 string, ServerType string, Data interface{}) bool {

	strServerType := ServerType
	for itr := MServer.Iterator(); itr.HasNext(); {
		k, v, _ := itr.Next()
		var key = ""
		var keyName = ""
		strsplit := Strings_Split(k.(string), "|") // key = serverid| Global_Server
		if len(strsplit) == 2 {
			for i := 0; i < len(strsplit); i++ {
				if i == 0 {
					key = strsplit[i]
				}
				// 获取链接的名字
				if i == len(strsplit)-1 {
					keyName = strsplit[i]
				}
				if key == StrMD5 && keyName == strServerType {
					// 发消息
					v.(interface{}).(*NetDataConn).PlayerSendMessage(Data)
					return true
				}
			}
		}
	}

	return false
}

// 发送给客户端
func (this *NetDataConn) SendClientDataFunc(StrMD5 string, ClientType string, Data interface{}) bool {

	fmt.Println("--------------------:", StrMD5)
	fmt.Println("--------------------:", ClientType)
	fmt.Println("--------------------:", Data)
	strClientType := ClientType
	for itr := M.Iterator(); itr.HasNext(); {
		k, v, _ := itr.Next()
		var key = ""
		var keyName = ""
		strsplit := Strings_Split(k.(string), "|") // key = serverid|connect
		if len(strsplit) == 2 {
			for i := 0; i < len(strsplit); i++ {
				if i == 0 {
					key = strsplit[i]
				}
				// 获取链接的名字
				if i == len(strsplit)-1 {
					keyName = strsplit[i]
				}
				if key == StrMD5 && keyName == strClientType {
					// 发消息
					v.(interface{}).(*NetDataConn).PlayerSendMessage(Data)
					return true
				}
			}
		}
	}

	return false
}

// 推送格式
// 例子参考
// www.521.mom
func (this *NetDataConn) XC_Data_Send_AllPlayer_State(StrMD5 string, Data interface{}) bool {
	//发给手机
	for itr := M.Iterator(); itr.HasNext(); {
		k, v, _ := itr.Next()
		var key = ""
		var keyName = ""
		// 拆分key
		strsplit := Strings_Split(k.(string), "|") // key = openid|XCN|name
		if len(strsplit) == 2 {
			// 拆分
			for i := 0; i < len(strsplit); i++ {
				if i == 0 {
					key = strsplit[i]
				}
				// 获取链接的名字
				if i == len(strsplit)-1 {
					keyName = strsplit[i]
				}
				if keyName == "connect" {
					// 发消息
					v.(interface{}).(*NetDataConn).PlayerSendMessage(Data)
				}
			}
		} else if len(strsplit) == 3 {
			// 拆分
			for i := 0; i < len(strsplit); i++ {
				if i == 1 {
					key = strsplit[i]
				}
				// 获取链接的名字
				if i == len(strsplit)-1 {
					keyName = strsplit[i]
				}
				if key == StrMD5 && keyName == "connect" {
					// 发消息
					v.(interface{}).(*NetDataConn).PlayerSendMessage(Data)
				}
			}
		}
	}

	return true
}
