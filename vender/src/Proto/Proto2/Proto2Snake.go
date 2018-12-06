package Proto2

import (
	"LollipopGo/LollipopGo/player"
)

const (
	ININSNAKE              = iota
	C2S_PlayerLoginSProto2 // PlayerLoginSProto2 == 1 登陆协议
	S2S_PlayerLoginSProto2 // S2S_PlayerLoginSProto2 == 2 登陆协议

	C2S_PlayerEntryGameProto2 // C2S_PlayerEntryGameProto2 == 3 进入游戏
	S2S_PlayerEntryGameProto2 // S2S_PlayerEntryGameProto2 == 4

	C2S_PlayerMoveProto2 // C2S_PlayerMoveProto2 == 5 移动操作
	S2S_PlayerMoveProto2 // S2S_PlayerMoveProto2 == 6

	C2S_PlayerAddGameProto2 // C2S_PlayerAddGameProto2 == 7 玩家进入匹配成功后进入游戏

)

//------------------------------------------------------------------------------

// C2S_PlayerAddGameProto2 玩家进入匹配成功后进入游戏
type C2S_PlayerAddGame struct {
	Protocol      int
	Protocol2     int
	OpenID        string // 玩家的唯一的标识, 另外一个玩家的唯一标识
	RoomID        int    // 房间ID
	PlayerHeadURL string // 玩家的头像数据
	Init_X        int    // 初始化X
	Init_Y        int    // 初始化Y
}

//------------------------------------------------------------------------------

// 移动操作
type C2S_PlayerMove struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 玩家的唯一的标识
	RoomID    int    // 房间ID
	OP_ULRDP  string // 玩家操作的方式：移动的方向
}

//  服务器广播给用户操作--同一个房间的，其他房间不广播
type S2S_PlayerMove struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 玩家的唯一的标识
	OP_ULRDP  string // 玩家操作的方式：移动的方向
}

//------------------------------------------------------------------------------

// 进入游戏匹配
type C2S_PlayerEntryGame struct {
	Protocol  int
	Protocol2 int
	Code      string //临时码
	Icode     int
}

//  返回数据操作
type S2S_PlayerEntryGame struct {
	Protocol  int
	Protocol2 int
	RoomID    int //房间ID
	Data      []int
	MapPlayer map[int]*player.PlayerSt // 玩家的结构信息
}

//------------------------------------------------------------------------------

// 登陆  客户端--> 服务器
type C2S_PlayerLoginS struct {
	Protocol   int
	Protocol2  int
	Login_Name string
	Login_PW   string
}

type S2S_PlayerLoginS struct {
	Protocol  int
	Protocol2 int
	Token     string // Token 的设计
}

//------------------------------------------------------------------------------
