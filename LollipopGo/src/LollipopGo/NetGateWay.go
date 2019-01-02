package main

import (
	"LollipopGo/LollipopGo/conf"
	"LollipopGo/LollipopGo/log"
	"LollipopGo/LollipopGo/match"
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
	case float64(Proto2.GW2G_PlayerMatchGameProto2):
		{
			// Global server 玩家匹配的协议
			this.GWPlayerMatchGameGL(ProtocolData)
		}
	default:
		panic("子协议：不存在！！！")
	}

	return
}

// Global server 返回给gateway server
func (this *NetDataConn) GWPlayerMatchGameGL(ProtocolData map[string]interface{}) {

	if ProtocolData["RoomUID"] == nil {
		log.Debug("Global server data is wrong:RoomUID is nil!")
		return
	}
	// 获取数据
	StrOpenID := ProtocolData["OpenID"].(string)
	StrRoomUID := ProtocolData["RoomUID"].(int)
	MatchPlayerST := ProtocolData["MatchPlayer"].(map[string]*match.RoomMatch)
	ChessBoard := ProtocolData["ChessBoard"].([4][4]int)
	iResultID := ProtocolData["ResultID"].(int)

	// 数据
	data_send := &Proto2.GW2G_PlayerMatchGame{
		Protocol:    Proto.G_GameGlobal_Proto, // 游戏主要协议
		Protocol2:   Proto2.GW2G_PlayerMatchGameProto2,
		OpenID:      StrOpenID, // 玩家唯一标识
		RoomUID:     StrRoomUID,
		MatchPlayer: MatchPlayerST,
		ChessBoard:  ChessBoard,
		ResultID:    iResultID,
	}

	// 发送数据  --
	// this.SendClientDataFunc(data_send.OpenID, "connect", data_send)
	// 发送给匹配的人的
	this.SendClientDataFunc(data_send.MatchPlayer[util.Int2str_LollipopGo(StrRoomUID)].PlayerAOpenID, "connect", data_send)
	this.SendClientDataFunc(data_send.MatchPlayer[util.Int2str_LollipopGo(StrRoomUID)].PlayerBOpenID, "connect", data_send)
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
	case float64(Proto2.C2GWS_PlayerChooseGameProto2):
		{
			// 功能函数处理 --  选择游戏列表的数据
			this.PlayerEntryGame(ProtocolData)
		}
	case float64(Proto2.C2GWS_PlayerChooseGameModeProto2):
		{
			// 功能函数处理 --  选择游戏对战类型
			this.PlayerChooseGameModeGame(ProtocolData)
		}
	default:
		panic("子协议：不存在！！！")
	}

	return
}

func (this *NetDataConn) PlayerChooseGameModeGame(ProtocolData map[string]interface{}) {
	if ProtocolData["OpenID"] == nil ||
		ProtocolData["RoomID"] == nil ||
		ProtocolData["Itype"] == nil {
		panic("选择游戏对战类型协议参数错误！")
		return
	}

	// 获取数据
	StrOpenID := ProtocolData["OpenID"].(string)
	iRoomID := ProtocolData["RoomID"].(int)
	Itype := ProtocolData["Itype"].(int)

	data := &Proto2.G2GW_PlayerMatchGame{
		Protocol:  Proto.G_GameGlobal_Proto,
		Protocol2: Proto2.G2GW_PlayerMatchGameProto2,
		OpenID:    StrOpenID, // 玩家唯一标识
		Itype:     Itype,     // Itype == 1：表示主动选择房间；Itype == 2：表示快速开始
		RoomID:    iRoomID,   // 房间ID
	}

	// 发送给global server
	this.SendServerDataFunc(strGlobalServer, "Global_Server", data)
	return
}

func (this *NetDataConn) PlayerEntryGame(ProtocolData map[string]interface{}) {
	if ProtocolData["OpenID"] == nil ||
		ProtocolData["GameID"] == nil {
		panic("进入游戏协议参数错误！")
		return
	}
	StrOpenID := ProtocolData["OpenID"].(string)
	StrGameID := ProtocolData["GameID"].(string)
	StrTimestamp := ProtocolData["Timestamp"].(string)
	_ = StrOpenID
	_ = StrTimestamp
	data := &Proto2.S2GWS_PlayerChooseGame{
		Protocol:  Proto.G_GateWay_Proto,
		Protocol2: Proto2.S2GWS_PlayerChooseGameProto2,
		RoomList:  conf.G_RoomList[StrGameID],
	}
	// 发送数据
	fmt.Println("房间列表:", data)
	this.PlayerSendMessage(data)
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
		UID:           StrPlayerUID,
		OpenID:        util.MD5_LollipopGO(StrPlayerUID + "GateWay"),
		PlayerName:    StrPlayerName,
		HeadUrl:       StrHeadUrl,
		Constellation: StrConstellation,
		PlayerSchool:  StrPlayerSchool,
		Sex:           StrSex,
		Token:         StrToken,
	}
	this.SendServerDataFunc(strGlobalServer, "Global_Server", data)

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
