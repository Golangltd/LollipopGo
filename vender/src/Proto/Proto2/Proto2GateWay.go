package Proto2

import (
	"LollipopGo/LollipopGo/player"
)

// G_GateWay_Proto
const (
	ININGATEWAY                      = iota // ININGATEWAY == 0
	C2GWS_PlayerLoginProto2                 // C2GWS_PlayerLoginProto2 == 1 登陆协议
	S2GWS_PlayerLoginProto2                 // S2GWS_PlayerLoginProto2 == 2
	GateWay_HeartBeatProto2                 // GateWay_HeartBeatProto2 == 3 心跳协议
	GateWay_RelinkProto2                    // GateWay_RelinkProto2 == 4 断线重新链接协议
	C2GWS_PlayerChooseGameProto2            // C2GWS_PlayerChooseGameProto2 == 5  // 玩家选择游戏
	S2GWS_PlayerChooseGameProto2            // S2GWS_PlayerChooseGameProto2 == 6
	C2GWS_PlayerChooseGameModeProto2        // C2GWS_PlayerChooseGameModeProto2 == 7  // 玩家选择游戏模式
	S2GWS_PlayerChooseGameModeProto2        // S2GWS_PlayerChooseGameModeProto2 == 8
	C2GWS_PlayerGameInitProto2              // C2GWS_PlayerGameInitProto2 == 9  // 匹配成功后，客户端下发获取初始化牌型
	S2GWS_PlayerGameInitProto2              // S2GWS_PlayerGameInitProto2 == 10
)

//------------------------------------------------------------------------------
// C2GWS_PlayerGameInitProto2
type C2GWS_PlayerGameInit struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	RoomUID   int
}

// S2GWS_PlayerGameInitProto2
type S2GWS_PlayerGameInit struct {
	Protocol   int
	Protocol2  int
	OpenID     string
	RoomUID    int
	ChessBoard [4][4]int // 棋盘的数据
}

//------------------------------------------------------------------------------
// C2GWS_PlayerChooseGameModeProto2
// 玩家选择游戏模式
type C2GWS_PlayerChooseGameMode struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 玩家唯一标识
	Itype     int    // Itype == 1：表示主动选择房间；Itype == 2：表示快速开始
	RoomID    int    // 房间ID
}

// S2GWS_PlayerChooseGameModeProto2
// 服务器返回数据
type S2GWS_PlayerChooseGameMode struct {
	Protocol    int
	Protocol2   int
	OpenID      string                 // 玩家唯一标识
	RoomUID     int                    // 房间ID；注意匹配失败或者超时，数据为空
	MatchPlayer map[string]interface{} // 匹配的玩家的信息；注意匹配失败或者超时，数据为空
	ChessBoard  [4][4]int              // 棋盘的数据
	ResultID    int                    // 结果ID
}

//------------------------------------------------------------------------------
// C2GWS_PlayerChooseGameProto2 玩家请求进入游戏
type C2GWS_PlayerChooseGame struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 玩家唯一的ID 信息
	GameID    string // 游戏ID
	Timestamp int    // 时间戳
}

// S2GWS_PlayerChooseGameProto2
type S2GWS_PlayerChooseGame struct {
	Protocol  int
	Protocol2 int
	RoomList  interface{}
}

//------------------------------------------------------------------------------
// 断线重连  网关
type GateWay_Relink struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	Timestamp int
}

//------------------------------------------------------------------------------

// C2GWS_PlayerLoginProto2
// 登陆  客户端--> 服务器
type C2GWS_PlayerLogin struct {
	Protocol      int
	Protocol2     int
	PlayerUID     string // APP 的UID
	PlayerName    string // 玩家的名字
	HeadUrl       string // 头像
	Constellation string // 星座
	PlayerSchool  string // 玩家的学校
	Sex           string // 性别
	Token         string
}

// S2GWS_PlayerLoginProto2
type S2GWS_PlayerLogin struct {
	Protocol      int
	Protocol2     int
	OpenID        string
	PlayerName    string                   // 玩家的名字
	HeadUrl       string                   // 头像
	Constellation string                   // 星座
	Sex           string                   // 性别
	GamePlayerNum map[string]interface{}   // 每个游戏的玩家的人数,global server获取
	RacePlayerNum map[string]interface{}   // 大奖赛列表
	Personal      map[string]interface{}   // 个人信息
	DefaultMsg    map[string]*player.MsgST // 默认跑马灯消息
	DefaultAward  map[string]interface{}   // 默认兑换列表
}

//------------------------------------------------------------------------------

// GateWay_HeartBeatProto2
// 心跳协议
type GateWay_HeartBeat struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 65位 玩家的唯一ID -- server ---> client (多数不需验证OpenID)
}
