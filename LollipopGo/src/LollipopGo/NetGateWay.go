package main

import (
	"LollipopGo/LollipopGo/log"
	"LollipopGo/LollipopGo/util"
	"Proto"
	"Proto/Proto2"
	"fmt"
)

//------------------------------------------------------------------------------
//------------------------------------------------------------------------------
//------------------------------------------------------------------------------
// Global Server 子协议的处理
func (this *NetDataConn) HandleCltProtocol2GL(protocol2 interface{}, ProtocolData map[string]interface{}) {

	switch protocol2 {
	case float64(Proto2.G2GW_ConnServerProto2):
		{
			// 网关主动链接进来，做数据链接的保存
			this.GLConnServerFunc(ProtocolData)
		}
	case float64(Proto2.GW2G_PlayerEntryHallProto2):
		{
			// Global server 返回给服务器
			this.GWPlayerLoginGL(ProtocolData)
		}
	default:
		panic("子协议：不存在！！！")
	}

	return
}

// Global server 返回给gateway server
func (this *NetDataConn) GWPlayerLoginGL(ProtocolData map[string]interface{}) {

	if ProtocolData["OpenID"] == nil {
		log.Debug("Global server data is wrong:OpenID is nil!")
		return
	}

	StrOpenID := ProtocolData["OpenID"].(string)
	StrPlayerName := ProtocolData["PlayerName"].(string)
	StrHeadUrl := ProtocolData["HeadUrl"].(string)
	StrConstellation := ProtocolData["Constellation"].(string)
	StrSex := ProtocolData["Sex"].(string)
	StGamePlayerNum := ProtocolData["GamePlayerNum"].(map[string]interface{})

	StRacePlayerNum := make(map[string]interface{})
	if ProtocolData["RacePlayerNum"] != nil {
		StRacePlayerNum = ProtocolData["RacePlayerNum"].(map[string]interface{})
	}
	StPersonal := ProtocolData["Personal"].(map[string]interface{})
	// StDefaultMsg := ProtocolData["DefaultMsg"].(map[string]*player.MsgST)
	// StDefaultMsg := ProtocolData["DefaultMsg"].(map[string]*player.MsgST)
	// StDefaultAward := ProtocolData["DefaultAward"].(map[string]interface{})

	// 发给客户端模拟
	data := &Proto2.S2GWS_PlayerLogin{
		Protocol:      6,
		Protocol2:     2,
		PlayerName:    StrPlayerName,
		HeadUrl:       StrHeadUrl,
		Constellation: StrConstellation,
		Sex:           StrSex,
		OpenID:        StrOpenID,
		GamePlayerNum: StGamePlayerNum,
		RacePlayerNum: StRacePlayerNum,
		Personal:      StPersonal,
		// DefaultMsg:    StDefaultMsg,
		// DefaultAward:  StDefaultAward,
	}
	// 发送数据  --
	this.SendClientDataFunc(data.OpenID, "connect", data)
	return
}

// Global server 保存
func (this *NetDataConn) GLConnServerFunc(ProtocolData map[string]interface{}) {
	if ProtocolData["ServerID"] == nil {
		panic("ServerID 数据为空!")
		return
	}

	fmt.Println("Global server conn entry gateway!!!")

	// Globla server 发过来的可以加密的数据
	StrServerID := ProtocolData["ServerID"].(string)
	strGlobalServer = StrServerID
	// 1 发送数据
	data := &Proto2.GW2G_ConnServer{
		Protocol:  9,
		Protocol2: 2,
		ServerID:  StrServerID,
	}
	// 发送数据
	this.PlayerSendMessage(data)

	// 2 保存Global的链接信息
	//================================推送消息处理===================================
	// 保存在线的玩家的数据信息
	onlineServer := &NetDataConn{
		Connection:    this.Connection, // 链接的数据信息
		MapSafeServer: this.MapSafeServer,
	}
	// 保存玩家数据到内存
	this.MapSafeServer.Put(StrServerID+"|Global_Server", onlineServer)
	//==============================================================================
	return
}

//------------------------------------------------------------------------------
//------------------------------------------------------------------------------
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

	StrPlayerUID := ProtocolData["PlayerUID"].(string)
	StrPlayerName := ProtocolData["PlayerName"].(string)
	StrHeadUrl := ProtocolData["HeadUrl"].(string)
	StrConstellation := ProtocolData["Constellation"].(string)
	StrPlayerSchool := ProtocolData["PlayerSchool"].(string)
	StrSex := ProtocolData["Sex"].(string)
	StrToken := ProtocolData["Token"].(string)

	// 1 将我们解析的数据 --> token --->  redis 验证等等
	// 主要看TTL的时间是否正确
	// 2 发送给Global server 获取数据  在线人数等
	data := &Proto2.G2GW_PlayerEntryHall{
		Protocol:      Proto.G_GameGlobal_Proto,
		Protocol2:     Proto2.G2GW_PlayerEntryHallProto2,
		OpenID:        util.MD5_LollipopGO(StrPlayerUID + "GateWay"),
		PlayerName:    StrPlayerName,
		HeadUrl:       StrHeadUrl,
		Constellation: StrConstellation,
		PlayerSchool:  StrPlayerSchool,
		Sex:           StrSex,
		Token:         StrToken,
	}
	this.SendServerDataFunc(strGlobalServer, "Global_Server", data)

	// 发给客户端模拟
	// data := &Proto2.S2GWS_PlayerLogin{
	// 	Protocol:  6,
	// 	Protocol2: 2,
	// 	OpenID:    util.MD5_LollipopGO(StrPlayerUID + "GateWay"),
	// }
	// 发送数据
	// this.PlayerSendMessage(data)
	// 保存玩家数据到内存 M
	//================================推送消息处理===================================
	// 保存在线的玩家的数据信息
	onlineUser := &NetDataConn{
		Connection: this.Connection, // 链接的数据信息
		MapSafe:    this.MapSafe,
	}
	// 保存玩家数据到内存
	this.MapSafe.Put(data.OpenID+"|connect", onlineUser)
	//==============================================================================

	return
}
