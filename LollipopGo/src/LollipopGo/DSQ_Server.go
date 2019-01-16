package main

import (
	"LollipopGo/LollipopGo/log"
	"Proto"
	"Proto/Proto2"
	"cache2go"
	"flag"
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
	"time"

	"LollipopGo/LollipopGo/util"
	_ "LollipopGo/ReadCSV"

	_ "LollipopGo/LollipopGo/player"

	"code.google.com/p/go.net/websocket"
)

var addrDSQ = flag.String("addrDSQ", "127.0.0.1:8888", "http service address") // 链接gateway
var ConnDSQ *websocket.Conn                                                    // 保存用户的链接信息，数据会在主动匹配成功后进行链接
var ConnDSQRPC *rpc.Client                                                     // 链接DB server
var DSQAllMap map[string]*RoomPlayerDSQ                                        // 游戏逻辑存储
var DSQ_qi = []int{                                                            // 1-8 A ;9-16 B ; 17 未翻牌; 18 已翻牌
	Proto2.Elephant, Proto2.Lion, Proto2.Tiger, Proto2.Leopard, Proto2.Wolf, Proto2.Dog, Proto2.Cat, Proto2.Mouse,
	Proto2.Mouse + Proto2.Elephant, Proto2.Mouse + Proto2.Lion, Proto2.Mouse + Proto2.Tiger, Proto2.Mouse + Proto2.Leopard,
	Proto2.Mouse + Proto2.Wolf, Proto2.Mouse + Proto2.Dog, Proto2.Mouse + Proto2.Cat, 2 * Proto2.Mouse}
var cacheDSQ *cache2go.CacheTable
var DSQLeftTime int = 30
var DSQTimeSpeed = time.Millisecond * 2000

type RoomPlayerDSQ struct {
	RoomID       int
	IsRoomID     bool
	OpenIDA      string
	OpenIDB      string
	OpenIDA_Seat int       // 默认是0； 1 表示坐下
	OpenIDB_Seat int       // 默认是0； 1 表示坐下
	Default      [4][4]int // 未翻牌的
	ChessData    [4][4]int // 棋盘数据
	ReChessData  [4][4]int // 重连棋盘数据
	AChessNum    int       // A的剩余的棋子的数量
	BChessNum    int       // B的剩余的棋子的数量
	GoAround     int       // 回合，如果每人出10次都没有吃子，系统推送平局;第七个回合提示数据 第10局平局
	LeftTime     int       // 剩余的时间
}

/*

	-----------------------------------------
	|										|
	|	[0,0]01	[1,0]02	[2,0]03	[3,0]04		|
	|										|
	|										|
	|	[0,1]05	[1,1]06	[2,1]07	[3,1]08		|
	|										|
	|										|
	|	[0,2]09	[1,2]10	[2,2]11	[3,2]12		|
	|									    |
	|										|
	|	[0,3]13	[1,3]14	[2,3]15	[3,3]16		|
	|										|
	-----------------------------------------

*/

// 初始化操作
func init() {
	//if strServerType == "DSQ" {
	if !initDSQGateWayNet() {
		fmt.Println("链接 gateway server 失败!")
		return
	}
	fmt.Println("链接 gateway server 成功!")
	// 初始化数据
	initDSQNetRPC()
	//}
	return
}

// 初始化RPC
func initDSQNetRPC() {
	client, err := jsonrpc.Dial("tcp", service)
	if err != nil {
		log.Debug("dial error:", err)
	}
	ConnDSQRPC = client
	cacheDSQ = cache2go.Cache("LollipopGo_DSQ")
}

// 初始化网关
func initDSQGateWayNet() bool {

	fmt.Println("用户客户端客户端模拟！")
	log.Debug("用户客户端客户端模拟！")
	url := "ws://" + *addrDSQ + "/GolangLtd"
	conn, err := websocket.Dial(url, "", "test://golang/")
	if err != nil {
		fmt.Println("err:", err.Error())
		return false
	}
	ConnDSQ = conn
	// 协程支持  --接受线程操作 全球协议操作
	go GameServerReceiveDSQ(ConnDSQ)
	// 发送链接的协议 ---》
	initConnDSQ(ConnDSQ)
	return true
}

// 链接到网关
func initConnDSQ(conn *websocket.Conn) {
	fmt.Println("---------------------------------")
	// 协议修改
	data := &Proto2.DSQ2GW_ConnServer{
		Protocol:  Proto.G_GameDSQ_Proto, // 游戏主要协议
		Protocol2: Proto2.DSQ2GW_ConnServerProto2,
		ServerID:  util.MD5_LollipopGO("8895" + "DSQ server"),
	}
	fmt.Println(data)
	// 2 发送数据到服务器
	PlayerSendToServer(conn, data)
	return
}

// 处理数据
func GameServerReceiveDSQ(ws *websocket.Conn) {
	for {
		var content string
		err := websocket.Message.Receive(ws, &content)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		// decode
		fmt.Println(strings.Trim("", "\""))
		fmt.Println(content)
		content = strings.Replace(content, "\"", "", -1)
		contentstr, errr := base64Decode([]byte(content))
		if errr != nil {
			fmt.Println(errr)
			continue
		}
		// 解析数据 --
		fmt.Println("返回数据：", string(contentstr))
		go SyncMeassgeFunDSQ(string(contentstr))
	}
}

// 链接分发 处理
func SyncMeassgeFunDSQ(content string) {
	var r Requestbody
	r.req = content

	if ProtocolData, err := r.Json2map(); err == nil {
		// 处理我们的函数
		HandleCltProtocolDSQ(ProtocolData["Protocol"], ProtocolData["Protocol2"], ProtocolData)
	} else {
		log.Debug("解析失败：", err.Error())
	}
}

//  主协议处理
func HandleCltProtocolDSQ(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{}) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			strerr := fmt.Sprintf("%s", err)
			//发消息给客户端
			ErrorST := Proto2.G_Error_All{
				Protocol:  Proto.G_Error_Proto,      // 主协议
				Protocol2: Proto2.G_Error_All_Proto, // 子协议
				ErrCode:   "80006",
				ErrMsg:    "亲，您发的数据的格式不对！" + strerr,
			}
			// 发送给玩家数据
			fmt.Println("Global server的主协议!!!", ErrorST)
		}
	}()

	// 协议处理
	switch protocol {
	case float64(Proto.G_GameDSQ_Proto):
		{ // DSQ Server 主要协议处理
			fmt.Println("DSQ server 主协议!!!")
			HandleCltProtocol2DSQ(protocol2, ProtocolData)
		}
	default:
		panic("主协议：不存在！！！")
	}
	return
}

// 子协议的处理
func HandleCltProtocol2DSQ(protocol2 interface{}, ProtocolData map[string]interface{}) {

	switch protocol2 {
	case float64(Proto2.GW2G_ConnServerProto2):
		{
			fmt.Println("gateway server DSQ server 数据信息！！！")
		}
	case float64(Proto2.GW2DSQ_InitGameProto2):
		{
			fmt.Println("网关请求获取棋盘初始化数据等")
			DSQ2GW_PlayerGameInitProto2Fucn(ConnDSQ, ProtocolData)
		}
	case float64(Proto2.GW2DSQ_PlayerStirChessProto2):
		{
			fmt.Println("玩家翻棋子的协议")
			GW2DSQ_PlayerStirChessProto2Fucn(ConnDSQ, ProtocolData)
		}
	case float64(Proto2.GW2DSQ_PlayerMoveChessProto2):
		{
			fmt.Println("玩家移动棋子的协议")
			GW2DSQ_PlayerMoveChessProto2Fucn(ConnDSQ, ProtocolData)
		}
	case float64(Proto2.GW2DSQ_PlayerGiveUpProto2):
		{
			fmt.Println("玩家放弃游戏的协议")
			GW2DSQ_PlayerGiveUpProto2Fucn(ConnDSQ, ProtocolData)
		}
	case float64(Proto2.GW2DSQ_PlayerRelinkGameProto2):
		{
			fmt.Println("玩家重新链接游戏的协议")
			GW2DSQ_PlayerRelinkGameProto2Fucn(ConnDSQ, ProtocolData)
		}
	default:
		panic("子协议：不存在！！！")
	}
	return
}

func GW2DSQ_PlayerRelinkGameProto2Fucn(conn *websocket.Conn, ProtocolData map[string]interface{}) {
	if ProtocolData["OpenID"] == nil ||
		ProtocolData["RoomUID"] == nil {
		panic(ProtocolData)
		return
	}
	StrOpenID := ProtocolData["OpenID"].(string)
	iRoomID := int(ProtocolData["RoomUID"].(float64))
	redata, lefttime := CacheGetReChessData(iRoomID)
	// 发送数据
	data := &Proto2.DSQ2GW_PlayerRelinkGame{
		Protocol:  Proto.G_GameDSQ_Proto,
		Protocol2: Proto2.DSQ2GW_PlayerRelinkGameProto2,
		OpenIDA:   StrOpenID,
		// OpenIDB   string
		LeftTime:  lefttime,
		ChessData: redata,
	}
	PlayerSendToServer(conn, data)
	return
}

func GW2DSQ_PlayerGiveUpProto2Fucn(conn *websocket.Conn, ProtocolData map[string]interface{}) {

	fmt.Println("玩家数据为空--GW2DSQ_PlayerGiveUpProto2Fucn")

	if ProtocolData["OpenID"] == nil ||
		ProtocolData["RoomUID"] == nil {
		panic(ProtocolData)
		return
	}
	StrOpenID := ProtocolData["OpenID"].(string)
	iRoomID := int(ProtocolData["RoomUID"].(float64))
	// 发送数据
	data := &Proto2.DSQ2GW_BroadCast_GameOver{
		Protocol:  Proto.G_GameDSQ_Proto,
		Protocol2: Proto2.DSQ2GW_BroadCast_GameOverProto2,
		IsDraw:    false,
		OpenIDA:   StrOpenID,                                          // 默认失败
		OpenIDB:   CacheGetOtherPlayerByPlayerUID(iRoomID, StrOpenID), // 胜利者的
		//		FailGameLev_Exp: "" + ",0",                                                                   // 格式: 1,10
		//		SuccGameLev_Exp: "" + ",10",                                                                  // 格式: 1,10
		//		FailPlayer:      nil,                                                                         // 失败者
		//		SuccPlayer:      nil,                                                                         // 胜利者
	}
	PlayerSendToServer(conn, data)
	// 删除房间数据
	cacheDSQ.Delete(iRoomID)

	return
	//--------------------------------------------------------------------------
	// 失败者*数据保存到db
	dbdata := &Proto2.DB_GameOver{
		OpenID:    data.OpenIDA,
		GameID:    Proto2.DSQ_GameID,
		GameLev:   0,
		GameExp:   0,
		GameScore: 0,
		GameItem:  "",
	}
	ModefyGamePlayerDataGameInfo(dbdata)
	// 胜利者*数据保存到db
	if true {
		dbdata := &Proto2.DB_GameOver{
			OpenID:    data.OpenIDB,
			GameID:    Proto2.DSQ_GameID,
			GameLev:   0,                  // 查询
			GameExp:   Proto2.DSQ_GameExp, // 增加的经验，等级DB/GW自己计算
			GameScore: 0,
			GameItem:  "",
		}
		ModefyGamePlayerDataGameInfo(dbdata)
	}
	//--------------------------------------------------------------------------
	return
	// gamelev := util.Sort_LollipopGo(conf.DSQGameExp, expdata+data.GameExp)
}

// 修改数据
func ModefyGamePlayerDataGameInfo(data *Proto2.DB_GameOver) interface{} {

	args := data
	// 返回的数据
	var reply bool
	//--------------------------------------------------------------------------
	// 同步调用
	// err = ConnRPC_GM.Call("Arith.SavePlayerDataST2DB", args, &reply)
	// if err != nil {
	// 	fmt.Println("Arith.SavePlayerDataST2DB call error:", err)
	// }
	// 异步调用
	divCall := ConnRPC_GM.Go("Arith.SavePlayerDataST2DB", args, &reply, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	fmt.Println(replyCall.Reply)
	//--------------------------------------------------------------------------
	// 返回的数据
	fmt.Println("ModefyGamePlayerDataGameInfo :", reply)
	return reply
}

func GW2DSQ_PlayerMoveChessProto2Fucn(conn *websocket.Conn, ProtocolData map[string]interface{}) {
	if ProtocolData["OpenID"] == nil ||
		ProtocolData["RoomUID"] == nil {
		panic(ProtocolData)
		return
	}

	StrOpenID := ProtocolData["OpenID"].(string)
	iRoomID := int(ProtocolData["RoomUID"].(float64))
	StrOldPos := ProtocolData["OldPos"].(string)
	iMoveDir := int(ProtocolData["MoveDir"].(float64))
	// if GetPlayerChupai(StrOpenID) {
	// 	data := &Proto2.DSQ2GW_PlayerMoveChess{
	// 		Protocol:  Proto.G_GameDSQ_Proto,
	// 		Protocol2: Proto2.DSQ2GW_PlayerMoveChessProto2,
	// 	}
	// 	data.ResultID = 60003
	// 	PlayerSendToServer(conn, data)
	// 	return
	// } else {
	// 	SetPlayerChupai(StrOpenID)
	// }

	// 1，是否可以移动（一定位置是都已经翻，移动的位置是否是自己人）
	// 2，移动成功，更新棋盘位置
	stropenida, stropenidb, strnewpos := CacheMoveChessIsUpdateData(iRoomID, StrOldPos, iMoveDir, StrOpenID)
	data := &Proto2.DSQ2GW_PlayerMoveChess{
		Protocol:  Proto.G_GameDSQ_Proto,
		Protocol2: Proto2.DSQ2GW_PlayerMoveChessProto2,
		// OpenIDA:   StrOpenID,
		// OpenIDB:   stropenidb,
		OpenIDA: stropenidb,
		OpenIDB: StrOpenID,
		RoomUID: iRoomID,
		OldPos:  StrOldPos,
		NewPos:  strnewpos,
	}
	// fmt.Println("++++++++++++++++++++", stropenidb)
	// fmt.Println("++++++++++++++++++++", stropenida)
	// fmt.Println("++++++++++++++++++++", data)
	if StrOpenID != stropenida {
		data.OpenIDB = stropenida
	}
	// 检测棋子
	CheckGameIsOver(iRoomID, StrOpenID)
	// fmt.Println("++++++++++++++++++++", data)
	PlayerSendToServer(conn, data)
	return
}

func GW2DSQ_PlayerStirChessProto2Fucn(conn *websocket.Conn, ProtocolData map[string]interface{}) {

	if ProtocolData["OpenID"] == nil ||
		ProtocolData["RoomUID"] == nil {
		panic(ProtocolData)
		return
	}

	StrOpenID := ProtocolData["OpenID"].(string)
	iRoomID := int(ProtocolData["RoomUID"].(float64))
	StrStirPos := ProtocolData["StirPos"].(string)

	data := &Proto2.DSQ2GW_PlayerStirChess{
		Protocol:  Proto.G_GameDSQ_Proto,
		Protocol2: Proto2.DSQ2GW_PlayerStirChessProto2,
		OpenID:    StrOpenID,
		OpenID_b:  "",
		StirPos:   StrStirPos,
		ResultID:  0,
	}
	// if GetPlayerChupai(StrOpenID) {
	// 	data.ResultID = 60003
	// 	PlayerSendToServer(conn, data)
	// 	return
	// } else {
	// 	SetPlayerChupai(StrOpenID)
	// }

	// 通过位置获取对应的数据
	_, idata := CacheGetChessDefaultData(iRoomID, StrStirPos, 2, 18)
	data.ChessNum = idata
	stropenid := CacheGetPlayerUID(iRoomID, StrOpenID)
	data.OpenID_b = stropenid
	// fmt.Println("++++++++++++++++++++", stropenid)
	// 发送数据
	PlayerSendToServer(conn, data)

	return
}

func DSQ2GW_PlayerGameInitProto2Fucn(conn *websocket.Conn, ProtocolData map[string]interface{}) {

	if ProtocolData["OpenID"] == nil ||
		ProtocolData["RoomID"] == nil {
		panic("玩家数据错误!!!")
		return
	}
	StrOpenID := ProtocolData["OpenID"].(string)
	StrRoomID := ProtocolData["RoomID"].(string)
	iRoomID := util.Str2int_LollipopGo(StrRoomID)
	retdata, bret, iret := CacheGetRoomDataByRoomID(iRoomID, StrOpenID)
	// iret := CacheGetRoomDataByPlayer(iRoomID, StrOpenID)
	if bret {
		data := &Proto2.DSQ2GW_InitGame{
			Protocol:  Proto.G_GameDSQ_Proto,
			Protocol2: Proto2.DSQ2GW_InitGameProto2,
			OpenID:    StrOpenID,
			RoomID:    StrRoomID,
			SeatNum:   iret,
			InitData:  retdata,
		}
		fmt.Println(data)
		PlayerSendToServer(conn, data)
		//
		//CacheSavePlayerUID(iRoomID, StrOpenID)
		return
	}
	data1 := [4][4]int{{2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1},
		{2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1},
		{2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1},
		{2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1, 2*Proto2.Mouse + 1}}
	var DSQ_qinei = []int{ // 1-8 A ;9-16 B ; 17 未翻牌; 18 已翻牌
		Proto2.Elephant, Proto2.Lion, Proto2.Tiger, Proto2.Leopard, Proto2.Wolf, Proto2.Dog, Proto2.Cat, Proto2.Mouse,
		Proto2.Mouse + Proto2.Elephant, Proto2.Mouse + Proto2.Lion, Proto2.Mouse + Proto2.Tiger, Proto2.Mouse + Proto2.Leopard,
		Proto2.Mouse + Proto2.Wolf, Proto2.Mouse + Proto2.Dog, Proto2.Mouse + Proto2.Cat, 2 * Proto2.Mouse}
	DSQ_Pai := InitDSQ(DSQ_qinei)
	savedata := &RoomPlayerDSQ{
		RoomID:      iRoomID,
		IsRoomID:    true,
		Default:     data1,
		ChessData:   DSQ_Pai,
		ReChessData: data1,
	}
	CacheSaveRoomData(iRoomID, savedata, StrOpenID)
	data := &Proto2.DSQ2GW_InitGame{
		Protocol:  Proto.G_GameDSQ_Proto,
		Protocol2: Proto2.DSQ2GW_InitGameProto2,
		OpenID:    StrOpenID,
		RoomID:    StrRoomID,
		SeatNum:   1,
		InitData:  DSQ_Pai,
	}
	PlayerSendToServer(conn, data)
	// cserli go CheckGameOfPlayerLeftTime(iRoomID, conn)
	return
}

//------------------------------------------------------------------------------
func CheckGameOfPlayerLeftTime(iRoomID int, conn *websocket.Conn) {

	for {
		select {
		case <-time.After(DSQTimeSpeed):
			{
				res, err1 := cacheDSQ.Value(iRoomID)
				if err1 != nil {
					continue
				}
				ilefttime := res.Data().(*RoomPlayerDSQ).LeftTime
				inowtime := int(util.GetNowUnix_LollipopGo())
				if inowtime-ilefttime > DSQLeftTime {
					data := &Proto2.DSQ2GW_BroadCast_GameOver{
						Protocol:  Proto.G_GameDSQ_Proto,
						Protocol2: Proto2.DSQ2GW_BroadCast_GameOverProto2,
						IsDraw:    false,
					}
					data.OpenIDA = res.Data().(*RoomPlayerDSQ).OpenIDA
					data.OpenIDB = res.Data().(*RoomPlayerDSQ).OpenIDB
					PlayerSendToServer(conn, data)
					return
				}
			}
		}
	}
	return
}

//------------------------------------------------------------------------------
// 吃掉棋子去检测
func CheckGameIsOver(iRoomID int, strpopenid string) bool {

	fmt.Println("吃掉棋子去检测")
	res, err1 := cacheDSQ.Value(iRoomID)
	if err1 != nil {
		// panic("没有对应数据")
		return false
	}

	// 结束
	data := &Proto2.DSQ2GW_BroadCast_GameOver{
		Protocol:  Proto.G_GameDSQ_Proto,
		Protocol2: Proto2.DSQ2GW_BroadCast_GameOverProto2,
		IsDraw:    false,
	}

	if res.Data().(*RoomPlayerDSQ).OpenIDA == strpopenid {
		idataB := res.Data().(*RoomPlayerDSQ).BChessNum
		fmt.Println("已经吃掉几个B", idataB)
		if idataB == 8 {
			data.OpenIDB = strpopenid
			data.OpenIDA = res.Data().(*RoomPlayerDSQ).OpenIDB
			PlayerSendToServer(ConnDSQ, data)
			return true
		}
		idataA := res.Data().(*RoomPlayerDSQ).AChessNum
		fmt.Println("已经吃掉几个A", idataA)
		if idataA == 8 {
			data.OpenIDB = res.Data().(*RoomPlayerDSQ).OpenIDB
			data.OpenIDA = strpopenid
			PlayerSendToServer(ConnDSQ, data)
			return true
		}
	}

	if res.Data().(*RoomPlayerDSQ).OpenIDB == strpopenid {
		idataA := res.Data().(*RoomPlayerDSQ).AChessNum
		fmt.Println("已经吃掉几个A", idataA)
		if idataA == 8 {
			data.OpenIDB = strpopenid
			data.OpenIDA = res.Data().(*RoomPlayerDSQ).OpenIDA
			PlayerSendToServer(ConnDSQ, data)
			return true
		}

		idataB := res.Data().(*RoomPlayerDSQ).BChessNum
		fmt.Println("已经吃掉几个B", idataB)
		if idataB == 8 {
			data.OpenIDB = res.Data().(*RoomPlayerDSQ).OpenIDA
			data.OpenIDA = strpopenid
			PlayerSendToServer(ConnDSQ, data)
			return true
		}
	}

	return false
}

//------------------------------------------------------------------------------
func SetRecord_A_B_chess(itype int) {
	// A B
	cacheDSQ.Add(itype, 0, 1)
	return
}

func GetRecord_A_B_chess(itype int) int {

	return 0
}

//------------------------------------------------------------------------------
func SetPlayerChupai(OpenID string) {
	cacheDSQ.Add(OpenID, 0, "exit")
}

func DelPlayerChupai(OpenID string) {
	cacheDSQ.Delete(OpenID)
}

func GetPlayerChupai(OpenID string) bool {
	ok := false
	_, err1 := cacheDSQ.Value(OpenID)
	if err1 == nil {
		ok = true
	}
	return ok
}

//------------------------------------------------------------------------------

func CacheSaveRoomData(iRoomID int, data *RoomPlayerDSQ, openid string) {
	fmt.Println("save data:", data)
	cacheDSQ.Add(iRoomID, 0, data)
	CacheSavePlayerUID(iRoomID, openid)
}

func CacheGetOtherPlayerByPlayerUID(iRoomID int, player string) string {
	res, err1 := cacheDSQ.Value(iRoomID)
	if err1 != nil {
		//panic("没有对应数据")
		return ""
	}
	if res.Data().(*RoomPlayerDSQ).OpenIDA == player {
		return res.Data().(*RoomPlayerDSQ).OpenIDB
	} else {
		return res.Data().(*RoomPlayerDSQ).OpenIDA
	}
	return ""
}

func CacheGetPlayerUID(iRoomID int, player string) string {
	res, err1 := cacheDSQ.Value(iRoomID)
	if err1 != nil {
		//panic("没有对应数据")
		return ""
	}
	fmt.Println("save data:", player)
	fmt.Println("save data:", res.Data().(*RoomPlayerDSQ).OpenIDA)
	fmt.Println("save data:", res.Data().(*RoomPlayerDSQ).OpenIDB)

	if res.Data().(*RoomPlayerDSQ).OpenIDA == player {
		return res.Data().(*RoomPlayerDSQ).OpenIDB
	}
	if res.Data().(*RoomPlayerDSQ).OpenIDB == player {
		return res.Data().(*RoomPlayerDSQ).OpenIDA
	}

	return ""
}

func CacheSavePlayerUID(iRoomID int, player string) {
	res, err1 := cacheDSQ.Value(iRoomID)
	if err1 != nil {
		//panic("没有对应数据")
		return
	}
	fmt.Println("result:", res.Data().(*RoomPlayerDSQ).OpenIDA)
	fmt.Println("result:", res.Data().(*RoomPlayerDSQ).OpenIDB)
	if len(res.Data().(*RoomPlayerDSQ).OpenIDA) == 0 {
		res.Data().(*RoomPlayerDSQ).OpenIDA = player
	} else {
		res.Data().(*RoomPlayerDSQ).OpenIDB = player
	}
	fmt.Println("result:", res.Data().(*RoomPlayerDSQ).OpenIDA)
	fmt.Println("result:", res.Data().(*RoomPlayerDSQ).OpenIDB)
	return
}

func CacheGetRoomDataByRoomID(iRoomID int, opneid string) ([4][4]int, bool, int) {
	res, err1 := cacheDSQ.Value(iRoomID)
	if err1 != nil {
		// panic("棋盘数据更新失败！")
		return [4][4]int{{}, {}, {}, {}}, false, 0
	}
	fmt.Println("n>1获取棋盘数据", iRoomID)
	if !res.Data().(*RoomPlayerDSQ).IsRoomID {
		return [4][4]int{{}, {}, {}, {}}, false, 0
	}

	if len(res.Data().(*RoomPlayerDSQ).OpenIDA) == 0 {
		res.Data().(*RoomPlayerDSQ).OpenIDA = opneid
		res.Data().(*RoomPlayerDSQ).OpenIDA_Seat = 1
		fmt.Println("----------------n>1获取棋盘数据res.Data().(*RoomPlayerDSQ).OpenIDB", res.Data().(*RoomPlayerDSQ).OpenIDB)
		fmt.Println("----------------n>1获取棋盘数据res.Data().(*RoomPlayerDSQ).OpenIDA", res.Data().(*RoomPlayerDSQ).OpenIDA)
		return res.Data().(*RoomPlayerDSQ).ChessData, true, 1
	}

	if len(res.Data().(*RoomPlayerDSQ).OpenIDB) == 0 {
		res.Data().(*RoomPlayerDSQ).OpenIDB = opneid
		res.Data().(*RoomPlayerDSQ).OpenIDB_Seat = 1
		fmt.Println("+++++++++++++n>1获取棋盘数据res.Data().(*RoomPlayerDSQ).OpenIDB", res.Data().(*RoomPlayerDSQ).OpenIDB)
		fmt.Println("++++++++++++++n>1获取棋盘数据res.Data().(*RoomPlayerDSQ).OpenIDA", res.Data().(*RoomPlayerDSQ).OpenIDA)
		return res.Data().(*RoomPlayerDSQ).ChessData, true, 0
	}

	return [4][4]int{{}, {}, {}, {}}, false, 0
}

func CacheGetRoomDataByPlayer(iRoomID int, opneid string) int {
	res, err1 := cacheDSQ.Value(iRoomID)
	if err1 != nil {
		//panic("棋盘数据更新失败！")
	}
	fmt.Println("n>1获取棋盘数据", iRoomID, opneid)

	if res.Data().(*RoomPlayerDSQ).OpenIDA == opneid {
		return res.Data().(*RoomPlayerDSQ).OpenIDA_Seat
	}

	if res.Data().(*RoomPlayerDSQ).OpenIDB == opneid {
		return res.Data().(*RoomPlayerDSQ).OpenIDB_Seat
	}

	if len(res.Data().(*RoomPlayerDSQ).OpenIDA) == 0 {
		res.Data().(*RoomPlayerDSQ).OpenIDA = opneid
		res.Data().(*RoomPlayerDSQ).OpenIDA_Seat = 1
		return 1
	}

	if len(res.Data().(*RoomPlayerDSQ).OpenIDB) == 0 {
		res.Data().(*RoomPlayerDSQ).OpenIDB = opneid
		res.Data().(*RoomPlayerDSQ).OpenIDB_Seat = 1
		return 0
	}
	fmt.Println("座位出错", iRoomID, opneid)
	return 0
}

func CacheUpdateRoomData(iRoomID int, Update_pos string, value int) {

	res, err1 := cacheDSQ.Value(iRoomID)
	if err1 != nil {
		//panic("棋盘数据更新失败！")
		return
	}

	ipos_x := 0
	ipos_y := 0
	strsplit := Strings_Split(Update_pos, ",")
	if len(strsplit) != 2 {
		panic("棋盘数据更新失败！")
		return
	}
	for i := 0; i < len(strsplit); i++ {
		if i == 0 {
			ipos_x = util.Str2int_LollipopGo(strsplit[i])
		} else {
			ipos_y = util.Str2int_LollipopGo(strsplit[i])
		}
	}
	fmt.Println("修改的棋盘的坐标", ipos_x, ipos_y)
	// 测试数据
	fmt.Println("result:", res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y])
	res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y] = value
	fmt.Println("result:", res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y])
	return
}

// 检查是否是七个回合没有吃棋子
func CheckIs7GoAround(iGoAround int, OpenIDA string, OpenIDB string, iRoomID int) {

	if iGoAround == 7 {
		fmt.Println("吃掉棋子去检测 --》7")
		data := &Proto2.DSQ_GameHint{
			Protocol:  Proto.G_GameDSQ_Proto,
			Protocol2: Proto2.DSQ_GameHintProto2,
			OpenIDA:   OpenIDA,
			OpenIDB:   OpenIDB,
		}
		PlayerSendToServer(ConnDSQ, data)
	} else if iGoAround == 10 {

		fmt.Println("吃掉棋子去检测 --》10")
		res, err1 := cacheDSQ.Value(iRoomID)
		if err1 != nil {
			return
		}

		// 结束
		data := &Proto2.DSQ2GW_BroadCast_GameOver{
			Protocol:  Proto.G_GameDSQ_Proto,
			Protocol2: Proto2.DSQ2GW_BroadCast_GameOverProto2,
			IsDraw:    true,
			OpenIDA:   OpenIDA,
			OpenIDB:   OpenIDB,
		}
		PlayerSendToServer(ConnDSQ, data)
		res.Data().(*RoomPlayerDSQ).GoAround = 0

		// 清除棋盘数据
		cacheDSQ.Delete(iRoomID)
	}
	return
}

// 获取重新链接数据
func CacheGetReChessData(iRoomID int) ([4][4]int, int) {

	res, err1 := cacheDSQ.Value(iRoomID)
	if err1 != nil {
		panic("棋盘数据获取数据失败！")
	}
	return res.Data().(*RoomPlayerDSQ).ReChessData, res.Data().(*RoomPlayerDSQ).LeftTime
}

// 移动期盼是否可以移动
func CacheMoveChessIsUpdateData(iRoomID int, Update_pos string, MoveDir int, stropenid string) (string, string, string) {
	res, err1 := cacheDSQ.Value(iRoomID)
	if err1 != nil {
		//panic("棋盘数据获取数据失败！")
		return "", "", ""
	}
	ipos_x := 0
	ipos_y := 0
	strsplit := Strings_Split(Update_pos, ",")
	if len(strsplit) != 2 {
		panic("棋盘数据获取数据失败！")
		return "", "", ""
	}
	for i := 0; i < len(strsplit); i++ {
		if i == 0 {
			ipos_x = util.Str2int_LollipopGo(strsplit[i])
		} else {
			ipos_y = util.Str2int_LollipopGo(strsplit[i])
		}
	}

	fmt.Println("修改的棋盘的坐标", ipos_x, ipos_y)
	iyunalaiX, iyunalaiY := ipos_x, ipos_y
	// 原来的值：
	iyuanlai := res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y]

	// 方向
	if MoveDir == Proto2.UP {
		ipos_y -= 1
	} else if MoveDir == Proto2.DOWN {
		ipos_y += 1
	} else if MoveDir == Proto2.LEFT {
		ipos_x -= 1
	} else if MoveDir == Proto2.RIGHT {
		ipos_x += 1
	}
	sendopenid, otheropenid := "", ""
	if res.Data().(*RoomPlayerDSQ).OpenIDA == stropenid {
		sendopenid = res.Data().(*RoomPlayerDSQ).OpenIDA
		otheropenid = res.Data().(*RoomPlayerDSQ).OpenIDB
	} else if res.Data().(*RoomPlayerDSQ).OpenIDB == stropenid {
		sendopenid = res.Data().(*RoomPlayerDSQ).OpenIDB
		otheropenid = res.Data().(*RoomPlayerDSQ).OpenIDA
	}
	ihoulai := res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y]
	strnewpos := util.Int2str_LollipopGo(ipos_x) + "," + util.Int2str_LollipopGo(ipos_y)
	if ihoulai == 0 {
		res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y] = iyuanlai
		res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = 0

		res.Data().(*RoomPlayerDSQ).ReChessData[ipos_x][ipos_y] = iyuanlai
		res.Data().(*RoomPlayerDSQ).ReChessData[iyunalaiX][iyunalaiY] = 0

		res.Data().(*RoomPlayerDSQ).GoAround++
		CheckIs7GoAround(res.Data().(*RoomPlayerDSQ).GoAround, sendopenid, otheropenid, iRoomID)

		res.Data().(*RoomPlayerDSQ).LeftTime = int(util.GetNowUnix_LollipopGo())
		return sendopenid, otheropenid, strnewpos
	}
	// 移动的位置
	bret, _ := CacheGetChessDefaultData(iRoomID, strnewpos, 1, 0)
	if !bret {
		res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y] = ihoulai
		res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = iyuanlai

		res.Data().(*RoomPlayerDSQ).ReChessData[ipos_x][ipos_y] = ihoulai
		res.Data().(*RoomPlayerDSQ).ReChessData[iyunalaiX][iyunalaiY] = iyuanlai

		res.Data().(*RoomPlayerDSQ).LeftTime = int(util.GetNowUnix_LollipopGo())
		return sendopenid, otheropenid, ""
	}
	// 判断是否可以吃，1 大小； 2 是都是同一方
	if iyuanlai > 8 && ihoulai > 8 {
		fmt.Println("都大于8", ipos_x, ipos_y)
		// res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y] = ihoulai
		// res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = iyuanlai
		return sendopenid, otheropenid, ""
	} else if iyuanlai <= 8 && ihoulai <= 8 {
		fmt.Println("都小于8", ipos_x, ipos_y)
		// res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y] = ihoulai
		// res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = iyuanlai
		return sendopenid, otheropenid, ""
	}

	if iyuanlai <= 8 && ihoulai > 8 {
		if iyuanlai == 1 && ihoulai-Proto2.Mouse == 8 {
			// 大象吃不了老鼠
			return sendopenid, otheropenid, ""
		}
		res.Data().(*RoomPlayerDSQ).GoAround = 0
		res.Data().(*RoomPlayerDSQ).LeftTime = int(util.GetNowUnix_LollipopGo())
		if iyuanlai == 8 && ihoulai-Proto2.Mouse == 1 {
			res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y] = iyuanlai
			res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = 0

			res.Data().(*RoomPlayerDSQ).ReChessData[ipos_x][ipos_y] = iyuanlai
			res.Data().(*RoomPlayerDSQ).ReChessData[iyunalaiX][iyunalaiY] = 0

			if stropenid == res.Data().(*RoomPlayerDSQ).OpenIDA {
				res.Data().(*RoomPlayerDSQ).BChessNum++
			} else {
				res.Data().(*RoomPlayerDSQ).AChessNum++
			}

			// CheckIs7GoAround(res.Data().(*RoomPlayerDSQ).GoAround, sendopenid, otheropenid, iRoomID)

			return sendopenid, otheropenid, strnewpos
		}
		if iyuanlai < ihoulai-Proto2.Mouse { // 可以吃
			res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y] = iyuanlai
			res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = 0

			res.Data().(*RoomPlayerDSQ).ReChessData[ipos_x][ipos_y] = iyuanlai
			res.Data().(*RoomPlayerDSQ).ReChessData[iyunalaiX][iyunalaiY] = 0
			if stropenid == res.Data().(*RoomPlayerDSQ).OpenIDA {
				res.Data().(*RoomPlayerDSQ).BChessNum++
			} else {
				res.Data().(*RoomPlayerDSQ).AChessNum++
			}
			return sendopenid, otheropenid, strnewpos

		} else if iyuanlai == ihoulai-Proto2.Mouse { // 同归于尽
			res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y] = 0
			res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = 0

			res.Data().(*RoomPlayerDSQ).ReChessData[ipos_x][ipos_y] = 0
			res.Data().(*RoomPlayerDSQ).ReChessData[iyunalaiX][iyunalaiY] = 0
			res.Data().(*RoomPlayerDSQ).AChessNum++
			res.Data().(*RoomPlayerDSQ).BChessNum++
			return sendopenid, otheropenid, strnewpos

		} else if iyuanlai > ihoulai-Proto2.Mouse { // 自毁
			res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = 0

			res.Data().(*RoomPlayerDSQ).ReChessData[iyunalaiX][iyunalaiY] = 0

			if stropenid == res.Data().(*RoomPlayerDSQ).OpenIDA {
				res.Data().(*RoomPlayerDSQ).BChessNum++
			} else {
				res.Data().(*RoomPlayerDSQ).AChessNum++
			}
			return sendopenid, otheropenid, strnewpos
		}

		return sendopenid, otheropenid, strnewpos

	} else if iyuanlai > 8 && ihoulai <= 8 {

		if iyuanlai-Proto2.Mouse == 1 && ihoulai == Proto2.Mouse {
			// 大象吃不了老鼠
			return "", "", ""
		}

		// sendopenid, otheropenid := "", ""
		// if res.Data().(*RoomPlayerDSQ).OpenIDA == stropenid {
		// 	sendopenid = res.Data().(*RoomPlayerDSQ).OpenIDA
		// 	otheropenid = res.Data().(*RoomPlayerDSQ).OpenIDB
		// } else if res.Data().(*RoomPlayerDSQ).OpenIDB == stropenid {
		// 	sendopenid = res.Data().(*RoomPlayerDSQ).OpenIDB
		// 	otheropenid = res.Data().(*RoomPlayerDSQ).OpenIDA
		// }
		res.Data().(*RoomPlayerDSQ).GoAround = 0
		res.Data().(*RoomPlayerDSQ).LeftTime = int(util.GetNowUnix_LollipopGo())
		if iyuanlai-Proto2.Mouse == 8 && ihoulai == 1 {
			res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y] = iyuanlai
			res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = 0

			res.Data().(*RoomPlayerDSQ).ReChessData[ipos_x][ipos_y] = iyuanlai
			res.Data().(*RoomPlayerDSQ).ReChessData[iyunalaiX][iyunalaiY] = 0
			if stropenid == res.Data().(*RoomPlayerDSQ).OpenIDA {
				res.Data().(*RoomPlayerDSQ).BChessNum++
			} else {
				res.Data().(*RoomPlayerDSQ).AChessNum++
			}
			return sendopenid, otheropenid, strnewpos
		}

		if iyuanlai-Proto2.Mouse < ihoulai { // 可以吃
			res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y] = iyuanlai
			res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = 0

			res.Data().(*RoomPlayerDSQ).ReChessData[ipos_x][ipos_y] = iyuanlai
			res.Data().(*RoomPlayerDSQ).ReChessData[iyunalaiX][iyunalaiY] = 0
			if stropenid == res.Data().(*RoomPlayerDSQ).OpenIDA {
				res.Data().(*RoomPlayerDSQ).BChessNum++
			} else {
				res.Data().(*RoomPlayerDSQ).AChessNum++
			}
			return sendopenid, otheropenid, strnewpos

		} else if iyuanlai-Proto2.Mouse == ihoulai { // 同归于尽
			res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y] = 0
			res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = 0

			res.Data().(*RoomPlayerDSQ).ReChessData[ipos_x][ipos_y] = 0
			res.Data().(*RoomPlayerDSQ).ReChessData[iyunalaiX][iyunalaiY] = 0
			res.Data().(*RoomPlayerDSQ).AChessNum++
			res.Data().(*RoomPlayerDSQ).BChessNum++
			return sendopenid, otheropenid, strnewpos

		} else if iyuanlai-Proto2.Mouse < ihoulai { // 自毁
			res.Data().(*RoomPlayerDSQ).ChessData[iyunalaiX][iyunalaiY] = 0

			res.Data().(*RoomPlayerDSQ).ReChessData[iyunalaiX][iyunalaiY] = 0

			if stropenid == res.Data().(*RoomPlayerDSQ).OpenIDA {
				res.Data().(*RoomPlayerDSQ).BChessNum++
			} else {
				res.Data().(*RoomPlayerDSQ).AChessNum++
			}
			return sendopenid, otheropenid, strnewpos
		}

		return sendopenid, otheropenid, strnewpos
	}

	return "", "", ""
}

// 获取默认棋牌数据是否翻开
// true 表示翻开了
// itype ==1 查询是否翻开
// itype ==2 修改数据
func CacheGetChessDefaultData(iRoomID int, Update_pos string, itype int, valve int) (bool, int) {
	res, err1 := cacheDSQ.Value(iRoomID)
	if err1 != nil {
		//panic("棋盘数据获取数据失败！")
		return false, -1
	}
	ipos_x := 0
	ipos_y := 0
	strsplit := Strings_Split(Update_pos, ",")
	if len(strsplit) != 2 {
		panic("棋盘数据获取数据失败！")
		return false, -1
	}
	for i := 0; i < len(strsplit); i++ {
		if i == 0 {
			ipos_x = util.Str2int_LollipopGo(strsplit[i])
		} else {
			ipos_y = util.Str2int_LollipopGo(strsplit[i])
		}
	}
	fmt.Println("修改的棋盘的坐标", ipos_x, ipos_y)
	if itype == 1 {
		// 获取
		fmt.Println("result:", res.Data().(*RoomPlayerDSQ).Default[ipos_x][ipos_y])
		idata := res.Data().(*RoomPlayerDSQ).Default[ipos_x][ipos_y]
		if idata == (2*Proto2.Mouse + 1) {
			return false, -1
		} else {
			return true, -1
		}
	} else if itype == 2 {
		// 修改翻盘结构
		fmt.Println("result:", res.Data().(*RoomPlayerDSQ).Default[ipos_x][ipos_y])
		res.Data().(*RoomPlayerDSQ).Default[ipos_x][ipos_y] = valve
		res.Data().(*RoomPlayerDSQ).LeftTime = int(util.GetNowUnix_LollipopGo())
		return true, res.Data().(*RoomPlayerDSQ).ChessData[ipos_x][ipos_y]
	}
	return false, -1
}

//------------------------------------------------------------------------------
// 初始化牌型
func InitDSQ(data1 []int) [4][4]int {

	data, erdata, j, k := data1, [4][4]int{}, 0, 0

	for i := 0; i < Proto2.Mouse*2; i++ {
		icount := util.RandInterval_LollipopGo(0, int32(len(data))-1)
		//fmt.Println("随机数：", icount)
		if len(data) == 1 {
			erdata[3][3] = data[0]
		} else {
			//------------------------------------------------------------------
			if int(icount) < len(data) {
				erdata[j][k] = data[icount]
				k++
				if k%4 == 0 {
					j++
					k = 0
				}
				data = append(data[:icount], data[icount+1:]...)
			} else {
				erdata[j][k] = data[icount]
				k++
				if k%4 == 0 {
					j++
					k = 0
				}
				data = data[:icount-1]
			}
			//------------------------------------------------------------------
		}
		//fmt.Println("生成的数据", erdata)
	}

	return erdata
}

// 判断棋子大小
// 玩家操作移动,操作协议
func CheckIsEat(fangx int, qizi int, qipan [4][4]int) (bool, int, [4][4]int) {
	if qizi > 16 || qizi < 1 {
		log.Debug("玩家发送棋子数据不对！")
		return false, Proto2.DATANOEXIT, qipan
	}
	// 1 寻找 玩家的棋子在棋牌的位置/或者这个棋子是否存在
	bret, Posx, posy := CheckChessIsExit(qizi, qipan)
	if bret {
		bret, iret, data := CheckArea(fangx, Posx, posy, qipan)
		return bret, iret, data
	} else {
		log.Debug("玩家发送棋子不存在！疑似外挂。")
		return false, Proto2.DATAERROR, qipan
	}

	return true, 0, qipan
}

// 检查棋盘中是不是存在
func CheckChessIsExit(qizi int, qipan [4][4]int) (bool, int, int) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if qipan[i][j] == qizi {
				return true, i, j
			}
		}
	}
	return false, 0, 0
}

// 边界判断
func CheckArea(fangx, iposx, iposy int, qipan [4][4]int) (bool, int, [4][4]int) {

	if fangx == Proto2.UP {
		if iposy == 0 {
			return false, Proto2.MOVEFAIL, qipan // 无法移动
		}
		data_yidong := qipan[iposx][iposy-1]
		data := qipan[iposx][iposy]
		// 判断是空地不
		if data_yidong == 0 {
			data_ret := UpdateChessData(Proto2.MOVE, iposx, iposy-1, iposx, iposy, qipan)
			return true, Proto2.MOVESUCC, data_ret // 移动成功
		}
		// 对比棋子大小
		if data < 9 {
			if data_yidong < 9 {
				return false, Proto2.TEAMMATE, qipan // 自己人
			} else {
				if data_yidong > data {
					data_ret := UpdateChessData(Proto2.DISAPPEAR, iposx, iposy-1, iposx, iposy, qipan)
					return true, Proto2.DISAPPEAR, data_ret // 自残
				} else if data_yidong == data {
					data_ret := UpdateChessData(Proto2.ALLDISAPPEAR, iposx, iposy-1, iposx, iposy, qipan)
					return true, Proto2.ALLDISAPPEAR, data_ret // 同归于尽
				} else if data_yidong < data {
					data_ret := UpdateChessData(Proto2.BEAT, iposx, iposy-1, iposx, iposy, qipan)
					return true, Proto2.BEAT, data_ret // 吃掉对方
				}
			}
		} else {
			if data_yidong > 9 {
				return false, Proto2.TEAMMATE, qipan // 自己人
			} else {
				if data_yidong > data {
					data_ret := UpdateChessData(Proto2.DISAPPEAR, iposx, iposy-1, iposx, iposy, qipan)
					return true, Proto2.DISAPPEAR, data_ret // 自残
				} else if data_yidong == data {
					data_ret := UpdateChessData(Proto2.ALLDISAPPEAR, iposx, iposy-1, iposx, iposy, qipan)
					return true, Proto2.ALLDISAPPEAR, data_ret // 同归于尽
				} else if data_yidong < data {
					data_ret := UpdateChessData(Proto2.BEAT, iposx, iposy-1, iposx, iposy, qipan)
					return true, Proto2.BEAT, data_ret // 吃掉对方
				}
			}
		}

	} else if fangx == Proto2.DOWN {

		if iposy == 3 {
			return false, Proto2.MOVEFAIL, qipan // 无法移动
		}
		data_yidong := qipan[iposx][iposy+1]
		data := qipan[iposx][iposy]
		// 判断是空地不
		if data_yidong == 0 {
			data_ret := UpdateChessData(Proto2.MOVE, iposx, iposy+1, iposx, iposy, qipan)
			return true, Proto2.MOVESUCC, data_ret // 移动成功
		}
		// 对比棋子大小
		if data < 9 {
			if data_yidong < 9 {
				return false, Proto2.TEAMMATE, qipan // 自己人
			} else {
				if data_yidong > data {
					data_ret := UpdateChessData(Proto2.DISAPPEAR, iposx, iposy+1, iposx, iposy, qipan)
					return true, Proto2.DISAPPEAR, data_ret // 自残
				} else if data_yidong == data {
					data_ret := UpdateChessData(Proto2.ALLDISAPPEAR, iposx, iposy+1, iposx, iposy, qipan)
					return true, Proto2.ALLDISAPPEAR, data_ret // 同归于尽
				} else if data_yidong < data {
					data_ret := UpdateChessData(Proto2.BEAT, iposx, iposy+1, iposx, iposy, qipan)
					return true, Proto2.BEAT, data_ret // 吃掉对方
				}
			}
		} else {
			if data_yidong > 9 {
				return false, Proto2.TEAMMATE, qipan // 自己人
			} else {
				if data_yidong > data {
					data_ret := UpdateChessData(Proto2.DISAPPEAR, iposx, iposy+1, iposx, iposy, qipan)
					return true, Proto2.DISAPPEAR, data_ret // 自残
				} else if data_yidong == data {
					data_ret := UpdateChessData(Proto2.ALLDISAPPEAR, iposx, iposy+1, iposx, iposy, qipan)
					return true, Proto2.ALLDISAPPEAR, data_ret // 同归于尽
				} else if data_yidong < data {
					data_ret := UpdateChessData(Proto2.BEAT, iposx, iposy+1, iposx, iposy, qipan)
					return true, Proto2.BEAT, data_ret // 吃掉对方
				}
			}
		}

	} else if fangx == Proto2.LEFT {
		if iposx == 0 {
			return false, Proto2.MOVEFAIL, qipan // 无法移动
		}
		data_yidong := qipan[iposx-1][iposy]
		data := qipan[iposx][iposy]
		// 判断是空地不
		if data_yidong == 0 {
			data_ret := UpdateChessData(Proto2.MOVE, iposx-1, iposy, iposx, iposy, qipan)
			return true, Proto2.MOVESUCC, data_ret // 移动成功
		}
		// 对比棋子大小
		if data < 9 {
			if data_yidong < 9 {
				return false, Proto2.TEAMMATE, qipan // 自己人
			} else {
				if data_yidong > data {
					data_ret := UpdateChessData(Proto2.DISAPPEAR, iposx-1, iposy, iposx, iposy, qipan)
					return true, Proto2.DISAPPEAR, data_ret // 自残
				} else if data_yidong == data {
					data_ret := UpdateChessData(Proto2.ALLDISAPPEAR, iposx-1, iposy, iposx, iposy, qipan)
					return true, Proto2.ALLDISAPPEAR, data_ret // 同归于尽
				} else if data_yidong < data {
					data_ret := UpdateChessData(Proto2.BEAT, iposx-1, iposy, iposx, iposy, qipan)
					return true, Proto2.BEAT, data_ret // 吃掉对方
				}
			}
		} else {
			if data_yidong > 9 {
				return false, Proto2.TEAMMATE, qipan // 自己人
			} else {
				if data_yidong > data {
					data_ret := UpdateChessData(Proto2.DISAPPEAR, iposx-1, iposy, iposx, iposy, qipan)
					return true, Proto2.DISAPPEAR, data_ret // 自残
				} else if data_yidong == data {
					data_ret := UpdateChessData(Proto2.ALLDISAPPEAR, iposx-1, iposy, iposx, iposy, qipan)
					return true, Proto2.ALLDISAPPEAR, data_ret // 同归于尽
				} else if data_yidong < data {
					data_ret := UpdateChessData(Proto2.BEAT, iposx-1, iposy, iposx, iposy, qipan)
					return true, Proto2.BEAT, data_ret // 吃掉对方
				}
			}
		}

	} else if fangx == Proto2.RIGHT {
		if iposx == 3 {
			return false, Proto2.MOVEFAIL, qipan // 无法移动
		}
		data_yidong := qipan[iposx+1][iposy]
		data := qipan[iposx][iposy]
		// 判断是空地不
		if data_yidong == 0 {
			data_ret := UpdateChessData(Proto2.MOVE, iposx+1, iposy, iposx, iposy, qipan)
			return true, Proto2.MOVESUCC, data_ret // 移动成功
		}
		// 对比棋子大小
		if data < 9 {
			if data_yidong < 9 {
				return false, Proto2.TEAMMATE, qipan // 自己人
			} else {
				if data_yidong > data {
					data_ret := UpdateChessData(Proto2.DISAPPEAR, iposx+1, iposy, iposx, iposy, qipan)
					return true, Proto2.DISAPPEAR, data_ret // 自残
				} else if data_yidong == data {
					data_ret := UpdateChessData(Proto2.ALLDISAPPEAR, iposx+1, iposy, iposx, iposy, qipan)
					return true, Proto2.ALLDISAPPEAR, data_ret // 同归于尽
				} else if data_yidong < data {
					data_ret := UpdateChessData(Proto2.BEAT, iposx+1, iposy, iposx, iposy, qipan)
					return true, Proto2.BEAT, data_ret // 吃掉对方
				}
			}
		} else {
			if data_yidong > 9 {
				return false, Proto2.TEAMMATE, qipan // 自己人
			} else {
				if data_yidong > data {
					data_ret := UpdateChessData(Proto2.DISAPPEAR, iposx+1, iposy, iposx, iposy, qipan)
					return true, Proto2.DISAPPEAR, data_ret // 自残
				} else if data_yidong == data {
					data_ret := UpdateChessData(Proto2.ALLDISAPPEAR, iposx+1, iposy, iposx, iposy, qipan)
					return true, Proto2.ALLDISAPPEAR, data_ret // 同归于尽
				} else if data_yidong < data {
					data_ret := UpdateChessData(Proto2.BEAT, iposx+1, iposy, iposx, iposy, qipan)
					return true, Proto2.BEAT, data_ret // 吃掉对方
				}
			}
		}
	}

	return false, Proto2.ITYPEINIY, qipan // 不存在的操作
}

// 更新棋盘数据
// 注：移动成功后，原来位置如何变化
//    目标的位置如何变化
//    fangxPos fangyPos表示变化的位置
//    iPosx iPosy 原来棋子的位置
func UpdateChessData(iType, fangxPos, fangyPos, iPosx, iPosy int, qipan [4][4]int) [4][4]int {
	// 保存副本数据
	data := qipan
	// 原来棋子的数据
	yd_data := data[iPosx][iPosy]
	//  更新棋盘 数据
	if iType == Proto2.MOVE { // 更新空地 ,1 更新原来棋盘的位置为0， 目标的位置为现在数据
		data[iPosx][iPosy] = 0
		data[fangxPos][fangyPos] = yd_data
	} else if iType == Proto2.DISAPPEAR { // 自残 ,1 更新原来棋盘的位置为0，目标的位置数据不变
		data[iPosx][iPosy] = 0
	} else if iType == Proto2.ALLDISAPPEAR { // 同归于尽 ,1 更新原来棋盘的位置为0，目标的位置数据为0
		data[iPosx][iPosy] = 0
		data[fangxPos][fangyPos] = 0
	} else if iType == Proto2.BEAT { // 吃掉对方 ,1 更新原来位置为0，目标位置为现在数据
		data[iPosx][iPosy] = 0
		data[fangxPos][fangyPos] = yd_data
	}
	return data
}
