package main

import (
	"LollipopGo/LollipopGo/util"
	"Proto"
	"Proto/Proto2"
	"fmt"
)

// 子协议的处理
func (this *NetDataConn) HandleCltProtocol2Snake(protocol2 interface{}, ProtocolData map[string]interface{}) {

	switch protocol2 {
	case float64(Proto2.C2S_PlayerEntryGameProto2):
		{
			fmt.Println("贪吃蛇:玩家进入游戏的协议!!!")
			// 玩家进入游戏的协议
			this.EntryGameSnake(ProtocolData)
		}
	case float64(Proto2.C2S_PlayerLoginSProto2):
		{
			// 登录协议
			this.LoginGameSnake(ProtocolData)
		}

	default:
		panic("子协议：不存在！！！")
	}

	return
}

// 玩家 进入游戏的协议
func (this *NetDataConn) EntryGameSnake(ProtocolData map[string]interface{}) {
	// 进入游戏 进行匹配
	fmt.Println("玩家进行匹配!!!")
	// 1 匹配算法 --》匹配在线玩家的

	// 2 返回数据
	data := &Proto2.S2S_PlayerEntryGame{
		Protocol:  Proto.G_Snake_Proto,
		Protocol2: Proto2.S2S_PlayerEntryGameProto2,
		RoomID:    1,
		MapPlayer: nil,
	}
	// 发送数据给客户端了
	this.PlayerSendMessage(data)
	return
}

// 登录游戏协议
func (this *NetDataConn) LoginGameSnake(ProtocolData map[string]interface{}) {
	// 数据链接信息需要保存
	if ProtocolData["Login_Name"] == nil ||
		ProtocolData["Login_PW"] == nil {
		panic("玩家登录协议不存在, 字段为空！")
		return
	}
	StrLogin_Name := ProtocolData["Login_Name"].(string)
	// StrLogin_PW := ProtocolData["Login_PW"].(string)
	// 数据库验证
	// 1 获取到UID 信息
	// 服务器-->客户端
	// 2 获取玩家的信息 --> player data
	// 3 保存到内存中的数据
	data := &Proto2.S2S_PlayerLoginS{
		Protocol:  Proto.G_Snake_Proto,
		Protocol2: Proto2.S2S_PlayerLoginSProto2,
		Token:     "123456789",
	}
	// 发送数据给客户端了
	this.PlayerSendMessage(data)
	// 玩家信息保存到 内存中
	// MD5信息操作
	strRoom := "UID" // 玩家房间的ID信息：信息的组成主要是  游戏ID+房间ID信息 确定数据是唯一的
	// 数据保存操作
	// --> 操作的
	this.MapSafe.Put(StrLogin_Name+"|"+util.MD5_LollipopGO(strRoom)+"|connect", "")
	return
}
