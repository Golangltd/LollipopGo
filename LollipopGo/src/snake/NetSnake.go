package main

import (
	"Proto"
	"Proto/Proto2"
	_ "fmt"

	"code.google.com/p/go.net/websocket"
)

// 发送坐标数据 --==--
// 数据 X Y Z
func Send_XYZ_Data(conn *websocket.Conn, strOpenID string, strRoomID string, OP_ULRDP string) {
	// 1 组装  发送数据的协议
	data := &Proto2.C2S_PlayerMove{
		Protocol:  Proto.G_Snake_Proto, // 游戏主要协议
		Protocol2: Proto2.C2S_PlayerMoveProto2,
		OpenID:    strOpenID, // 随机生产的数据，时间戳
		RoomID:    1,
		OP_ULRDP:  OP_ULRDP,
	}
	// 2 发送数据到服务器
	PlayerSendToServer(conn, data)
}
