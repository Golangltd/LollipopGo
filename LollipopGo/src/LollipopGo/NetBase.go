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

	// 网络信息
	NetDataConntmp := &NetDataConn{
		Connection: ws,
		StrMd5:     "",
		MapSafe:    M,
	}
	// 指针接受者  处理消息
	NetDataConntmp.PullFromClient()
}

// 公用的send函数
func PlayerSendToServer(conn *websocket.Conn, data interface{}) {

	// 2 结构体转换成json数据
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
