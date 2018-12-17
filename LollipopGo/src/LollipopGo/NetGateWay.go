package main

import (
	"LollipopGo/LollipopGo/util"
	"Proto/Proto2"
)

// Global Server 子协议的处理
func (this *NetDataConn) HandleCltProtocol2GL(protocol2 interface{}, ProtocolData map[string]interface{}) {

	switch protocol2 {
	case float64(Proto2.G2GW_ConnServerProto2):
		{
			// 网关主动链接进来，做数据链接的保存
			this.GLConnServerFunc(ProtocolData)
		}
	default:
		panic("子协议：不存在！！！")
	}

	return
}

// Global server 保存
func (this *NetDataConn) GLConnServerFunc(ProtocolData map[string]interface{}) {
	if ProtocolData["ServerID"] == nil {
		panic("ServerID 数据为空!")
		return
	}

	// Globla server 发过来的可以加密的数据
	StrServerID := ProtocolData["ServerID"].(string)

	// 1 保存Global的链接信息
	// 2 发送数据

	return
}

//------------------------------------------------------------------------------

// client 子协议的处理
func (this *NetDataConn) HandleCltProtocol2GW(protocol2 interface{}, ProtocolData map[string]interface{}) {

	switch protocol2 {
	case float64(Proto2.C2GWS_PlayerLoginProto2):
		{
			// 功能函数处理 --  用户登陆协议
			this.GWPlayerLogin(ProtocolData)
		}
	case float64(Proto2.GateWay_HeartBeatProto2):
		{
			// 功能函数处理 --  心跳函数处理
			this.GWHeartBeat(ProtocolData)
		}
	default:
		panic("子协议：不存在！！！")
	}

	return
}

func (this *NetDataConn) GWHeartBeat(ProtocolData map[string]interface{}) {
	if ProtocolData["OpenID"] == nil {
		panic("心跳协议参数错误！")
		return
	}

	StrOpenID := ProtocolData["OpenID"].(string)
	// 将我们解析的数据 --> token --->  redis 验证等等
	// 主要看TTL的时间是否正确
	data := &Proto2.GateWay_HeartBeat{
		Protocol:  6,
		Protocol2: 3,
		OpenID:    StrOpenID,
	}
	// 发送数据
	this.PlayerSendMessage(data)
	return
}

func (this *NetDataConn) GWPlayerLogin(ProtocolData map[string]interface{}) {
	if ProtocolData["Token"] == nil ||
		ProtocolData["PlayerUID"] == nil {
		panic("网关登陆协议错误！！！")
		return
	}

	StrToken := ProtocolData["Token"].(string)
	StrPlayerUID := ProtocolData["PlayerUID"].(string)
	_ = StrToken
	// 将我们解析的数据 --> token --->  redis 验证等等
	// 主要看TTL的时间是否正确
	data := &Proto2.S2GWS_PlayerLogin{
		Protocol:  6,
		Protocol2: 2,
		OpenID:    util.MD5_LollipopGO(StrPlayerUID + "GateWay"),
	}
	// 发送数据
	this.PlayerSendMessage(data)

	return
}
