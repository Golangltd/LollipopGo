package main

import (
	"fmt"

	"LollipopGo/LollipopGo/util"

	"Proto"
	"Proto/Proto2"
	"encoding/base64"
	"encoding/json"

	"code.google.com/p/go.net/websocket"
)

// 公用的send函数
func PlayerSendToServer(conn *websocket.Conn, data interface{}) bool {

	// 2 结构体转换成json数据
	jsons, err := json.Marshal(data)
	if err != nil {
		fmt.Println("err:", err.Error())
		return false
	}

	errq := websocket.Message.Send(conn, jsons)
	if errq != nil {
		fmt.Println(errq)
		return false
	}
	return true
}

// 解码
func base64Decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

// 匹配服务器
func initMatch(conn *websocket.Conn) {

	// 1 组装  进入的协议
	data := &Proto2.C2S_PlayerEntryGame{
		Protocol:  Proto.G_Snake_Proto, // 游戏主要协议
		Protocol2: Proto2.C2S_PlayerEntryGameProto2,
		Code:      util.UTCTime_LollipopGO(), // 随机生产的数据，时间戳
	}
	// fmt.Println(data)
	// 2 发送数据到服务器
	PlayerSendToServer(conn, data)

}

// 登录服务器
// 登陆的用户名及密码都是随机暂时
func initLogin(conn *websocket.Conn) {
	// 组装数据
	data := &Proto2.C2S_PlayerLoginS{
		Protocol:   Proto.G_Snake_Proto, // 游戏主要协议
		Protocol2:  Proto2.C2S_PlayerLoginSProto2,
		Login_Name: util.UTCTime_LollipopGO(),
		Login_PW:   util.UTCTime_LollipopGO(),
	}
	// 发送数据
	PlayerSendToServer(conn, data)
}
