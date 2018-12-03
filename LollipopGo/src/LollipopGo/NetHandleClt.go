package main

import (
	"Proto"
	"Proto/Proto2"
	//"glog-master"
)

// 子协议的处理
func (this *NetDataConn) HandleCltProtocol2Net(protocol2 interface{}, ProtocolData map[string]interface{}) {

	switch protocol2 {
	case float64(Proto2.Net_HeartBeatProto2):
		{
			// 功能函数处理 --  心跳
			this.HeartBeat(ProtocolData)
		}
	case float64(Proto2.Net_RelinkProto2):
		{
			// 功能函数处理 --  重连
			this.Relink(ProtocolData)
		}
	default:
		panic("子协议：不存在！！！")
	}

	return
}

// 重新链接
func (this *NetDataConn) Relink(ProtocolData map[string]interface{}) {
	// 1 解析数据
	// 2 update 网络数据
	// 3 登陆流程？？  第二期 讲
	if ProtocolData["OpenID"] == nil {
		panic("心跳协议数据错误！！！")
		return
	}
	StrOpenID := ProtocolData["OpenID"].(string)
	_ = StrOpenID

	// 保存玩家数据
	playerdata := &NetDataConn{
		Connection: this.Connection,
	}

	this.MapSafe.Put("PlayerUID"+"|connect", playerdata)
	// 保存 --
	G_PlayerData["123456"] = playerdata

	// cache --- player --  check

	// 服务器-->客户端
	data := &Proto2.Net_Relink{
		Protocol:  Proto.GameNet_Proto,
		Protocol2: Proto2.Net_RelinkProto2,
		ISucc:     true,
	}
	// 发送数据给客户端了
	this.PlayerSendMessage(data)
	return
}

// 1sss --- 触发一次
func (this *NetDataConn) HeartBeat(ProtocolData map[string]interface{}) {
	// 1 解析 协议数据
	// 2 通过玩家的唯一ID 去保存心跳数据 map[]data
	// 3 timer -- 超时踢人

	//	glog.Info(ProtocolData["OpenID"])
	if ProtocolData["OpenID"] == nil {
		panic("心跳协议数据错误！！！")
		return
	}
	StrOpenID := ProtocolData["OpenID"].(string)
	if len(StrOpenID) == 65 {
		G_PlayerNet[StrOpenID]++
		// 防止溢出
		if G_PlayerNet[StrOpenID] > 100 {
			G_PlayerNet[StrOpenID] = 1
		}
		// 返回数据
		// 服务器-->客户端
		data := &Proto2.Net_HeartBeat{
			Protocol:  Proto.GameNet_Proto,
			Protocol2: Proto2.Net_HeartBeatProto2,
			OpenID:    StrOpenID,
		}
		// 发送数据给客户端了
		this.PlayerSendMessage(data)
	} else {
		panic("心跳协议数据错误！！！")
	}

	return
}
