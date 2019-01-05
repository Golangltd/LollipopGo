package Proto2

import (
	"LollipopGo/LollipopGo/match"
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
	C2GWS_QuitMatchProto2                   // C2GWS_QuitMatchProto2 == 11 退出协议
	S2GWS_QuitMatchProto2                   // S2GWS_QuitMatchProto2 == 12
	GateWay_LogoutProto2                    // GateWay_LogoutProto2  == 13 玩家登出

	/*
	   斗兽棋 0表示空 1-8 A方，9-16 B方
	*/
	C2GWS_PlayerStirChessProto2 // C2GWS_PlayerStirChessProto2 == 14   玩家翻棋子
	S2GWS_PlayerStirChessProto2 // S2GWS_PlayerStirChessProto2 == 15   广播同一个桌子上的,且接受到此协议后，已经移动的再无法移动棋子，对手获取操作权限
	C2GWS_PlayerMoveChessProto2 // C2GWS_PlayerMoveChessProto2 == 16   玩家移动
	S2GWS_PlayerMoveChessProto2 // S2GWS_PlayerMoveChessProto2 == 17   广播同一个桌子上的,且接受到此协议后，已经移动的再无法移动棋子，对手获取操作权限
	C2GWS_PlayerGiveUpProto2    // C2GWS_PlayerGiveUpProto2 == 18  玩家放弃、认输
	BroadCast_GameOverProto2    // BroadCast_GameOverProto2 == 19  广播玩家游戏结束
)

//------------------------------------------------------------------------------
// GateWay_RelinkProto2
// 重新链接网关
type GateWay_Relink struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 玩家唯一ID
	ResultID  int    // 结果ID == 1表示成功； 0：表示失败
}

//------------------------------------------------------------------------------
// C2GWS_QuitMatchProto2
// 玩家退出匹配
type C2GWS_QuitMatch struct {
	Protocol  int
	Protocol2 int
	OpenID    string
}

// S2GWS_QuitMatchProto2
// 玩家退出匹配  服务器清理数据
type S2GWS_QuitMatch struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	ResultID  int // 结果ID == 1表示成功； 0：表示失败
}

//------------------------------------------------------------------------------
// GateWay_LogoutProto2
// 玩家登出
type GateWay_Logout struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 玩家唯一ID
	ResultID  int    // 结果ID
}

//------------------------------------------------------------------------------
// BroadCast_GameOverProto2
type BroadCast_GameOver struct {
	Protocol        int
	Protocol2       int
	FailGameLev_Exp string                      // 格式: 1,10
	SuccGameLev_Exp string                      // 格式: 1,10
	FailPlayer      map[string]*player.PlayerSt // 失败者
	SuccPlayer      map[string]*player.PlayerSt // 胜利者
}

//------------------------------------------------------------------------------
// C2GWS_PlayerGiveUpProto2
// 认输
type C2GWS_PlayerGiveUp struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	RoomUID   int
}

//------------------------------------------------------------------------------
// C2GWS_PlayerMoveChessProto2
type C2GWS_PlayerMoveChess struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	RoomUID   int
	OldPos    string // 原来坐标
	MoveDir   int    // 移动的方向，UP == 1，DOWN 	== 2，LEFT 	== 3，RIGHT 	== 4
}

// S2GWS_PlayerMoveChessProto2
// 广播 同一个房间
type S2GWS_PlayerMoveChess struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	RoomUID   int
	OldPos    string // 原来坐标
	NewPos    string // 新坐标
}

//------------------------------------------------------------------------------
// C2GWS_PlayerStirChessProto2
type C2GWS_PlayerStirChess struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	RoomUID   int
	StirPos   string // 翻动的位置 格式: x,y
}

// S2GWS_PlayerStirChessProto2
// 广播
type S2GWS_PlayerStirChess struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 谁翻动了棋子
	StirPos   string // 翻动的位置  格式:x,y
	ChessNum  int    // 1 - 16 正数
}

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
	ChessBoard [4][4]int // 棋盘的数据 0表示空
}

//------------------------------------------------------------------------------
// C2GWS_PlayerChooseGameModeProto2
// 玩家选择游戏模式
type C2GWS_PlayerChooseGameMode struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 玩家唯一标识
	Itype     string // Itype == 1：表示主动选择房间；Itype == 2：表示快速开始
	RoomID    string // 房间ID
}

// S2GWS_PlayerChooseGameModeProto2
// 服务器返回数据
type S2GWS_PlayerChooseGameMode struct {
	Protocol    int
	Protocol2   int
	OpenID      string                      // 玩家唯一标识
	RoomUID     int                         // 房间ID；注意匹配失败或者超时，数据为空
	MatchPlayer map[string]*match.RoomMatch // 匹配的玩家的信息；注意匹配失败或者超时，数据为空
	ChessBoard  []interface{}               // 棋盘的数据
	ResultID    int                         // 结果ID
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
type GateWay_Relink1 struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	Timestamp string
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
