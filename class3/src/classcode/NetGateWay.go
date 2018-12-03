package main

import (
	"Proto/Proto2"
)

func init() {

	return
}

// 子协议的处理
func (this *NetDataConn) HandleCltProtocol2GW(protocol2 interface{}, ProtocolData map[string]interface{}) {

	switch protocol2 {
	case float64(Proto2.C2GWS_PlayerLoginProto2):
		{
			// 功能函数处理 --  用户登陆协议
			this.GWPlayerLogin(ProtocolData)
		}
	default:
		panic("子协议：不存在！！！")
	}

	return
}

// 用户奔跑的协议
func (this *NetDataConn) GWPlayerLogin(ProtocolData map[string]interface{}) {
	if ProtocolData["Token"] == nil {
		panic("网关登陆协议错误！！！")
		return
	}

	StrToken := ProtocolData["Token"].(string)
	_ = StrToken
	// 将我们解析的数据 --> token --->  redis 验证等等
	return
}
