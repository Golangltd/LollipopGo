package Proto2

import (
	"LollipopGo/LollipopGo/player"
)

const (
	ININGATEWAY                 = iota // ININGATEWAY == 0
	C2GWS_PlayerLoginProto2            // C2GWS_PlayerLoginProto2 == 1 登陆协议
	S2GWS_PlayerLoginProto2            // S2GWS_PlayerLoginProto2 == 2
	GateWay_HeartBeatProto2            // GateWay_HeartBeatProto2 == 3 心跳协议
	GateWay_RelinkProto2               // GateWay_RelinkProto2 == 4 断线重新链接协议
	C2GWS_PlayerEntryGameProto2        // C2GWS_PlayerEntryGameProto2 == 5 玩家请求进入游戏
	S2GWS_PlayerEntryGameProto2        // S2GWS_PlayerEntryGameProto2 == 6
)

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
// C2GWS_PlayerEntryGameProto2 玩家请求进入游戏
type C2GWS_PlayerEntryGame struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 玩家唯一的ID 信息
	GameID    string // 游戏ID
	Timestamp int    // 时间戳
}

// S2GWS_PlayerEntryGameProto2
type S2GWS_PlayerEntryGame struct {
	Protocol  int
	Protocol2 int
}

//------------------------------------------------------------------------------
// 断线重连  网关
type GateWay_Relink struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	Timestamp int // 时间戳
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
	Sex           string // 性别
	Token         string
}

// S2GWS_PlayerLoginProto2
type S2GWS_PlayerLogin struct {
	Protocol      int
	Protocol2     int
	OpenID        string
	GamePlayerNum map[string]interface{}   // 每个游戏的玩家的人数,global server获取
	RacePlayerNum map[string]interface{}   // 大奖赛列表
	Personal      *player.PlayerSt         // 个人信息
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
