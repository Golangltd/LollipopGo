package main

import (
	"Proto"
	"Proto/Proto2"
	"time"

	"code.google.com/p/go.net/websocket"
)

// 定时发心跳
func Timer(conn *websocket.Conn) {
	itimer := time.NewTicker(50)
	for {
		select {
		case <-itimer.C:

			strOpenID := "A123456789A123456789A123456789A123456789A123456789A123456789A1234"
			// 执行的代码
			//			data := &Proto2.Net_HeartBeat{
			//				Protocol:  Proto.GameNet_Proto,
			//				Protocol2: Proto2.Net_HeartBeatProto2,
			//				OpenID:    strOpenID,
			//			}
			//发送
			//PlayerSendToServer(conn, data)
			// -----------------------------------------------------------------
			// run的消息
			datapro := &Proto2.C2S_PlayerRun{
				Protocol:  Proto.GameData_Proto,
				Protocol2: Proto2.C2S_PlayerRunProto2,
				OpenID:    strOpenID,
				StrRunX:   "22",
				StrRunY:   "22",
				StrRunZ:   "22",
			}
			PlayerSendToServer(conn, datapro)
		}
	}
}
