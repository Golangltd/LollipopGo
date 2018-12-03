package main

import (
	"encoding/json"
	"fmt"

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
