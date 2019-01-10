package main

import (
	"LollipopGo/LollipopGo/conf"
	"LollipopGo/LollipopGo/error"
	"LollipopGo/LollipopGo/log"
	"LollipopGo/LollipopGo/match"
	"Proto"
	"Proto/Proto2"
	"flag"
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
	"time"

	"LollipopGo/LollipopGo/util"
	"LollipopGo/ReadCSV"

	"LollipopGo/LollipopGo/player"

	"code.google.com/p/go.net/websocket"
)

/*
  匹配、活动服务器
	1 匹配玩家活动
*/

var addrG = flag.String("addrG", "127.0.0.1:8888", "http service address")
var Conn *websocket.Conn
var ConnRPC *rpc.Client

func init() {
	if !initGateWayNet() {
		fmt.Println("链接 gateway server 失败!")
		return
	}
	fmt.Println("链接 gateway server 成功!")
	initNetRPC()

	return
}

func initNetRPC() {
	client, err := jsonrpc.Dial("tcp", service)
	if err != nil {
		log.Debug("dial error:", err)
		//panic("dial RPC Servre error")
		return
	}
	ConnRPC = client
}

func initGateWayNet() bool {
	fmt.Println("用户客户端客户端模拟！")
	url := "ws://" + *addrG + "/GolangLtd"
	conn, err := websocket.Dial(url, "", "test://golang/")
	if err != nil {
		fmt.Println("err:", err.Error())
		return false
	}
	Conn = conn
	go GameServerReceiveG(Conn)
	go TimeMsgNotice(Conn)
	initConn(Conn)
	return true
}

// 处理数据
func GameServerReceiveG(ws *websocket.Conn) {
	for {
		var content string
		err := websocket.Message.Receive(ws, &content)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(strings.Trim("", "\""))
		fmt.Println(content)
		content = strings.Replace(content, "\"", "", -1)
		contentstr, errr := base64Decode([]byte(content))
		if errr != nil {
			fmt.Println(errr)
			continue
		}
		go SyncMeassgeFunG(string(contentstr))
	}
}

// 链接分发 处理
func SyncMeassgeFunG(content string) {
	var r Requestbody
	r.req = content

	if ProtocolData, err := r.Json2map(); err == nil {
		HandleCltProtocolG(ProtocolData["Protocol"], ProtocolData["Protocol2"], ProtocolData)
	} else {
		log.Debug("解析失败：", err.Error())
	}
}

//  主协议处理
func HandleCltProtocolG(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{}) {
	// defer func() { // 必须要先声明defer，否则不能捕获到panic异常
	// 	if err := recover(); err != nil {
	// 		strerr := fmt.Sprintf("%s", err)
	// 		//发消息给客户端
	// 		ErrorST := Proto2.G_Error_All{
	// 			Protocol:  Proto.G_Error_Proto,      // 主协议
	// 			Protocol2: Proto2.G_Error_All_Proto, // 子协议
	// 			ErrCode:   "80006",
	// 			ErrMsg:    "亲，您发的数据的格式不对！" + strerr,
	// 		}
	// 		// 发送给玩家数据
	// 		fmt.Println("Global server的主协议!!!", ErrorST)
	// 	}
	// }()

	// 协议处理
	switch protocol {
	case float64(Proto.G_GameGlobal_Proto):
		{ // Global Server 主要协议处理
			fmt.Println("Global server 主协议!!!")
			HandleCltProtocol2Glogbal(protocol2, ProtocolData)

		}
	default:
		panic("主协议：不存在！！！")
	}
	return
}

// 子协议的处理
func HandleCltProtocol2Glogbal(protocol2 interface{}, ProtocolData map[string]interface{}) {

	switch protocol2 {
	case float64(Proto2.GW2G_ConnServerProto2):
		{ // 网关返回数据
			fmt.Println("gateway server 返回给global server 数据信息！！！")
		}
	case float64(Proto2.G2GW_PlayerEntryHallProto2):
		{
			G2GW_PlayerEntryHallProto2Fucn(Conn, ProtocolData)
		}
	case float64(Proto2.G2GW_PlayerMatchGameProto2):
		{
			fmt.Println("玩家请求玩家匹配！")
			G2GW_PlayerMatchGameProto2Fucn(Conn, ProtocolData)
		}
	case float64(Proto2.GW2G_PlayerQuitMatchGameProto2):
		{
			fmt.Println("玩家主动退出匹配！")
			G2GW_PlayerQuitMatchGameProto2Fucn(Conn, ProtocolData)
		}
	case float64(Proto2.GW2G_GetPlayerEmailListProto2):
		{
			fmt.Println("获取玩家邮件列表！")
			G2GW_PlayerGetPlayerEmailListProto2Fucn(Conn, ProtocolData)
		}
	case float64(Proto2.GW2G_ReadOrDelPlayerEmailProto2):
		{
			fmt.Println("玩家邮件列表读取！")
			G2GW_PlayerReadOrDelPlayerEmailProto2Fucn(Conn, ProtocolData)
		}
	default:
		panic("子协议：不存在！！！")
	}
	return
}

func G2GW_PlayerReadOrDelPlayerEmailProto2Fucn(conn *websocket.Conn, ProtocolData map[string]interface{}) {
	if ProtocolData["OpenID"] == nil {
		panic("读取邮件错误!")
	}
	StrOpenID := ProtocolData["OpenID"].(string)
	iItype := int(ProtocolData["Itype"].(float64))
	iEmailID := int(ProtocolData["EmailID"].(float64))

	data_send := &Proto2.G2GW_ReadOrDelPlayerEmail{
		Protocol:  Proto.G_GameGlobal_Proto,
		Protocol2: Proto2.G2GW_ReadOrDelPlayerEmailProto2,
		OpenID:    StrOpenID,
		Itype:     iItype,
		EmailID:   iEmailID,
	}
	// 1:读取打开，2：删除，3：领取附件
	if iItype == 1 {
		data_send.Itype = 0
		if EmailDatatmp[iEmailID] != nil {
			EmailDatatmp[iEmailID].IsOpen = true
			data_send.Itype = 1
		}
	} else if iItype == 2 {
		delete(EmailDatatmp, iEmailID)
	} else if iItype == 3 {
		EmailDatatmp[iEmailID].IsGet = true
	}

	PlayerSendToServer(conn, data_send)
	return
}

//------------------------------------------------------------------------------

var EmailDatatmp map[int]*player.EmailST
var ItemListtmp map[int]*player.ItemST
var PaoMaDeng map[int]*player.MsgST
var iicounmsg int = 3
var iicounemail int = 6

func init() {
	EmailDatatmp = make(map[int]*player.EmailST)
	ItemListtmp = make(map[int]*player.ItemST)
	PaoMaDeng = make(map[int]*player.MsgST)

	if true {
		data := new(player.EmailST)
		data.ID = 1
		data.Name = "测试邮件1"
		data.Sender = "test1"
		data.Type = 1
		data.Time = int(util.GetNowUnix_LollipopGo())
		data.Content = "测试邮件内容1"
		data.IsAdd_ons = false
		data.IsOpen = false
		data.IsGet = false
		EmailDatatmp[data.ID] = data
	}

	if true {
		data := new(player.EmailST)
		data.ID = 2
		data.Name = "测试邮件2"
		data.Sender = "test2"
		data.Type = 4
		data.Time = int(util.GetNowUnix_LollipopGo())
		data.Content = "测试邮件内容2"
		data.IsAdd_ons = false
		data.IsOpen = false
		data.IsGet = false
		EmailDatatmp[data.ID] = data
	}

	if true {
		data := new(player.EmailST)
		data.ID = 3
		data.Name = "测试邮件3"
		data.Sender = "test3"
		data.Type = 1
		data.Time = int(util.GetNowUnix_LollipopGo())
		data.Content = "测试邮件内容3"
		data.IsAdd_ons = true
		data.IsOpen = false
		data.IsGet = false

		if true {
			dataitem := new(player.ItemST)
			dataitem.ID = 1
			dataitem.Icon = ""
			dataitem.Name = "M卡"
			dataitem.Itype = 1
			dataitem.Num = 10
			ItemListtmp[dataitem.ID] = dataitem
		}

		data.ItemList = ItemListtmp
		EmailDatatmp[data.ID] = data
	}

	if true {
		data := new(player.EmailST)
		data.ID = 4
		data.Name = "测试邮件4"
		data.Sender = "test4"
		data.Type = 4
		data.Time = int(util.GetNowUnix_LollipopGo())
		data.Content = "测试邮件内容1"
		data.IsAdd_ons = false
		data.IsOpen = false
		data.IsGet = true
		EmailDatatmp[data.ID] = data
	}

	if true {
		data := new(player.EmailST)
		data.ID = 5
		data.Name = "测试邮件5"
		data.Sender = "test5"
		data.Type = 1
		data.Time = int(util.GetNowUnix_LollipopGo())
		data.Content = "测试邮件内容1"
		data.IsAdd_ons = false
		data.IsOpen = true
		data.IsGet = true
		EmailDatatmp[data.ID] = data
	}

	if true {
		data := new(player.EmailST)
		data.ID = 6
		data.Name = "测试邮件6"
		data.Sender = "test6"
		data.Type = 4
		data.Time = int(util.GetNowUnix_LollipopGo())
		data.Content = "测试邮件内容3"
		data.IsAdd_ons = true
		data.IsOpen = true
		data.IsGet = true

		if true {
			dataitem := new(player.ItemST)
			dataitem.ID = 1
			dataitem.Icon = ""
			dataitem.Name = "M卡"
			dataitem.Itype = 1
			dataitem.Num = 10
			ItemListtmp[dataitem.ID] = dataitem
		}

		data.ItemList = ItemListtmp
		EmailDatatmp[data.ID] = data
	}
	//--------------------------------------------------------------------------
	// DefaultMsg    map[string]*player.MsgST    // 默认跑马灯消息
	if true {
		data := new(player.MsgST)
		data.MsgID = 1
		data.MsgType = player.MsgType1
		data.MsgDesc = "系统消息：充值998，送B站24K纯金哥斯拉"
		PaoMaDeng[data.MsgID] = data
	}
	if true {
		data := new(player.MsgST)
		data.MsgID = 2
		data.MsgType = player.MsgType2
		data.MsgDesc = "恭喜【XXX玩家】在XX比赛中获得xxx奖励"
		PaoMaDeng[data.MsgID] = data
	}
	if true {
		data := new(player.MsgST)
		data.MsgID = 3
		data.MsgType = player.MsgType3
		data.MsgDesc = "恭喜【XXX玩家】在兑换中心成功兑换SSS"
		PaoMaDeng[data.MsgID] = data
	}
	return
}

func TimeMsgNotice(conn *websocket.Conn) {

	// if GL_type != "8894" {
	// 	return
	// }
	fmt.Println("TimeMsgNotice")
	for {
		select {
		case <-time.After(time.Second * 10):
			{
				iicounmsg++
				iicounemail++
				MsgNoticeFuncbak(conn)
				EmailNoticeFunc(conn)
			}
		}
	}
}

func EmailNoticeFunc(conn *websocket.Conn) {
	EmailDatatmpbak := make(map[int]*player.EmailST)

	if true {
		data := new(player.EmailST)
		data.ID = iicounemail
		data.Name = "测试邮件5"
		data.Sender = "test5"
		data.Type = 1
		data.Time = int(util.GetNowUnix_LollipopGo())
		data.Content = "测试邮件内容1"
		data.IsAdd_ons = false
		data.IsOpen = true
		data.IsGet = true
		EmailDatatmpbak[data.ID] = data
	}

	data_send := &Proto2.G_Broadcast_NoticePlayerEmail{
		Protocol:  Proto.G_GameGlobal_Proto,
		Protocol2: Proto2.G_Broadcast_NoticePlayerEmailProto2,
		OpenID:    "6412121cbb2dc2cb9e460cfee7046be2",
		EmailData: EmailDatatmpbak,
	}

	fmt.Println("邮件通知:", data_send)
	PlayerSendToServer(conn, data_send)
	return
}

// 全服通知
func MsgNoticeFuncbak(conn *websocket.Conn) {
	PaoMaDengbak := make(map[int]*player.MsgST)
	if true {
		data := new(player.MsgST)
		data.MsgID = iicounmsg
		data.MsgType = player.MsgType1
		data.MsgDesc = "系统消息：充值998，送B站24K纯金哥斯拉"
		PaoMaDengbak[data.MsgID] = data
	}

	data_send := &Proto2.G_Broadcast_MsgNoticePlayer{
		Protocol:  Proto.G_GameGlobal_Proto,
		Protocol2: Proto2.G_Broadcast_MsgNoticePlayerProto2,
		OpenID:    "6412121cbb2dc2cb9e460cfee7046be2",
		MsgData:   PaoMaDengbak,
	}

	fmt.Println("消息通知:", data_send)
	PlayerSendToServer(conn, data_send)

	return
}

//------------------------------------------------------------------------------

func G2GW_PlayerGetPlayerEmailListProto2Fucn(conn *websocket.Conn, ProtocolData map[string]interface{}) {
	if ProtocolData["OpenID"] == nil {
		panic("获取玩家列表!")
	}
	StrOpenID := ProtocolData["OpenID"].(string)

	data_send := &Proto2.G2GW_GetPlayerEmailList{
		Protocol:  Proto.G_GameGlobal_Proto,
		Protocol2: Proto2.G2GW_GetPlayerEmailListProto2,
		OpenID:    StrOpenID,
		EmailData: EmailDatatmp,
	}
	PlayerSendToServer(conn, data_send)
	return
}

// 玩家主动退出匹配
func G2GW_PlayerQuitMatchGameProto2Fucn(conn *websocket.Conn, ProtocolData map[string]interface{}) {
	if ProtocolData["OpenID"] == nil {
		panic("玩家主动退出匹配!")
		return
	}
	StrOpenID := ProtocolData["OpenID"].(string)
	// 玩家主动退出
	match.SetQuitMatch(StrOpenID)
	// 发送消息
	data_send := &Proto2.G2GW_PlayerQuitMatchGame{
		Protocol:  Proto.G_GameGlobal_Proto,
		Protocol2: Proto2.G2GW_PlayerQuitMatchGameProto2,
		OpenID:    StrOpenID,
		ResultID:  0,
	}
	PlayerSendToServer(conn, data_send)
	return
}

// 玩家匹配
func G2GW_PlayerMatchGameProto2Fucn(conn *websocket.Conn, ProtocolData map[string]interface{}) {
	if ProtocolData["OpenID"] == nil ||
		ProtocolData["RoomID"] == nil ||
		ProtocolData["Itype"] == nil {
		panic("选择游戏对战类型协议参数错误！")
		return
	}
	StrOpenID := ProtocolData["OpenID"].(string)
	StrRoomID := ProtocolData["RoomID"].(string) //  匹配数据
	StrItype := ProtocolData["Itype"].(string)   //  1 是正常匹配 2 是快速匹配

	// 数据
	data_send := &Proto2.GW2G_PlayerMatchGame{
		Protocol:  Proto.G_GameGlobal_Proto, // 游戏主要协议
		Protocol2: Proto2.GW2G_PlayerMatchGameProto2,
		OpenID:    StrOpenID, // 玩家唯一标识
		// RoomUID:     0,
		// MatchPlayer: nil,
		// ChessBoard:  {{}, {}, {}, {}},
		ResultID: 0,
	}
	if match.GetMatchQueue(StrOpenID) {
		data_send.ResultID = Error.IsMatch
		PlayerSendToServer(conn, data_send)
		return
	}
	match.SetMatchQueue(StrOpenID)

	if StrItype == "2" { //快速匹配
		PlayerSendToServer(conn, data_send)
		return
	}

	data := conf.RoomListDatabak[StrRoomID]
	fmt.Println("针对某房间ID去获取，相应的数据的", conf.RoomListDatabak, data.NeedLev, StrRoomID)
	dataplayer := DB_Save_RoleSTBak(StrOpenID)
	match.Putdata(dataplayer)
	s := string([]byte(data.NeedLev)[2:])
	if util.Str2int_LollipopGo(s) > dataplayer.Lev {
		data_send.ResultID = Error.Lev_lack
		PlayerSendToServer(conn, data_send)
		return
	} else if util.Str2int_LollipopGo(data.NeedPiece) > dataplayer.CoinNum {
		data_send.ResultID = Error.Coin_lack
		PlayerSendToServer(conn, data_send)
		return
	}

	if len(match.MatchData) > 1 {
		dar := <-match.MatchData_Chan
		data_send.MatchPlayer = dar
		fmt.Println(data_send)
		PlayerSendToServer(conn, data_send)
		match.DelMatchQueue(StrOpenID)
	} else {
		go PlayerMatchTime(conn, StrOpenID, data_send)
	}

	return
}

func PlayerMatchTime(conn *websocket.Conn, OpenID string, data_send *Proto2.GW2G_PlayerMatchGame) {
	icount := 0
	for {
		select {
		case <-time.After(match.PlaterMatchSpeed):
			{
				fmt.Println(icount)
				if icount >= 30 {
					PlayerSendToServer(conn, data_send)
					return
				}

				if len(match.MatchData_Chan) > 1 {
					dar := <-match.MatchData_Chan
					data_send.MatchPlayer = dar
					fmt.Println(data_send)
					PlayerSendToServer(conn, data_send)
					match.DelMatchQueue(OpenID)
					return
				}
				icount++
			}
		}
	}
}

// 保存数据都DB 人物信息
func DB_Save_RoleSTBak(openid string) *player.PlayerSt {

	args := player.PlayerSt{
		OpenID: openid,
	}

	var reply *player.PlayerSt
	// 异步调用【结构的方法】
	if ConnRPC != nil {
		// ConnRPC.Call("Arith.GetPlayerST2DB", args, &reply) 同步调用
		divCall := ConnRPC.Go("Arith.GetPlayerST2DB", args, &reply, nil)
		replyCall := <-divCall.Done
		_ = replyCall.Reply
	} else {
		fmt.Println("ConnRPC == nil")
	}
	return reply
}

func G2GW_PlayerEntryHallProto2Fucn(conn *websocket.Conn, ProtocolData map[string]interface{}) {

	StrUID := ProtocolData["UID"].(string)
	StrOpenID := ProtocolData["OpenID"].(string)
	StrPlayerName := ProtocolData["PlayerName"].(string)
	StrHeadUrl := ProtocolData["HeadUrl"].(string)
	StrSex := ProtocolData["Sex"].(string)
	StrConstellation := ProtocolData["Constellation"].(string)
	StrPlayerSchool := ProtocolData["PlayerSchool"].(string)
	StrToken := ProtocolData["Token"].(string)
	_ = StrToken

	// 获取在线人数
	ddd := make(map[string]interface{})
	csv.M_CSV.LollipopGo_RLockRange(ddd)
	// 查询数据库,找出游戏服务器的uid信息
	// 返回的数据操作
	datadb := DB_Save_RoleST(StrUID, StrPlayerName, StrHeadUrl, StrPlayerSchool, StrSex, StrConstellation, 0, 0, 2000, 0, 0)
	fmt.Println("--------------------------:", datadb)
	// 个人数据
	personalmap := make(map[string]*player.PlayerSt)
	personalmap["1"] = &datadb
	_ = personalmap["1"].OpenID

	// 组装数据
	data := &Proto2.GW2G_PlayerEntryHall{
		Protocol:      Proto.G_GameGlobal_Proto, // 游戏主要协议
		Protocol2:     Proto2.GW2G_PlayerEntryHallProto2,
		OpenID:        StrOpenID,
		PlayerName:    StrPlayerName,
		HeadUrl:       StrHeadUrl,
		Constellation: StrConstellation,
		Sex:           StrSex,
		GamePlayerNum: ddd,
		RacePlayerNum: nil,
		Personal:      personalmap,
		DefaultMsg:    PaoMaDeng,
		DefaultAward:  nil,
		IsNewEmail:    true,
	}

	icount := 0
	for _, value := range EmailDatatmp {
		idata := value.IsOpen
		if idata {
			icount++
		}
		if icount == len(EmailDatatmp) {
			data.IsNewEmail = false
			break
		}
	}

	fmt.Println(data)
	PlayerSendToServer(conn, data)
	return

}

// 保存数据都DB 人物信息
func DB_Save_RoleST(uid, strname, HeadURL, StrPlayerSchool, Sex, Constellation string, Lev, HallExp, CoinNum, MasonryNum, MCard int) player.PlayerSt {

	args := player.PlayerSt{
		UID:           util.Str2int_LollipopGo(uid),
		VIP_Lev:       0,
		Name:          strname,
		HeadURL:       HeadURL,
		Sex:           Sex,
		PlayerSchool:  StrPlayerSchool,
		Lev:           Lev,
		HallExp:       HallExp,
		CoinNum:       CoinNum,
		MasonryNum:    MasonryNum,
		MCard:         MCard,
		Constellation: Constellation,
		OpenID:        util.MD5_LollipopGO(uid),
	}

	var reply player.PlayerSt
	// 异步调用【结构的方法】
	if ConnRPC != nil {
		// ConnRPC.Call("Arith.SavePlayerST2DB", args, &reply) 同步调用
		divCall := ConnRPC.Go("Arith.SavePlayerST2DB", args, &reply, nil)
		replyCall := <-divCall.Done
		_ = replyCall.Reply
	} else {
		fmt.Println("ConnRPC == nil")
	}
	return reply
}

func initConn(conn *websocket.Conn) {
	data := &Proto2.G2GW_ConnServer{
		Protocol:  Proto.G_GameGlobal_Proto,
		Protocol2: Proto2.G2GW_ConnServerProto2,
		ServerID:  util.MD5_LollipopGO("8894" + "Global server"),
	}
	PlayerSendToServer(conn, data)
	return
}
