package Proto2

import (
	_ "LollipopGo/LollipopGo/player"
)

// G_GameGlobal_Proto == 9  负责全局的游戏逻辑 的子协议
// 注：server类型为单点
const (
	ININGlobal            = iota // 0
	G2GW_ConnServerProto2        // G2GW_ConnServerProto2 == 1 Global主动链接 gateway
	GW2G_ConnServerProto2        // GW2G_ConnServerProto2 == 2 选择链接

	GW2G_HeartBeatProto2 // GW2G_HeartBeatProto2 == 3      心跳协议  保活的操作

	G2GW_PlayerMatchProto2 // G2GW_PlayerMatchProto2 == 4 玩家发送匹配的协议
	GW2G_PlayerMatchProto2 // GW2G_PlayerMatchProto2 == 5 服务器返回数据对应匹配机制
)

//------------------------------------------------------------------------------
// G2GW_PlayerMatchProto2  玩家发送匹配的协议
type G2GW_PlayerMatch struct {
	Protocol   int
	Protocol2  int
	GWServerID int    // 网关的ID信息，主要是数据统计需要
	Type       int    // 匹配类型，1 1V1 ,2 2V2 ,3 5V5
	OpenID     string // 用户唯一ID
}

// GW2G_PlayerMatchProto2 服务器返回数据对应匹配机制
type GW2G_PlayerMatch struct {
	Protocol  int
	Protocol2 int
	Type      int    // 匹配类型，1 1V1 ,2 2V2 ,3 5V5
	RoomInfo  string // 匹配成功后的房间中的信息
}

//------------------------------------------------------------------------------
// GW2G_HeartBeatProto2  心跳协议
type GW2G_HeartBeat struct {
	Protocol  int
	Protocol2 int
	ServerID  int //全局配置 唯一的也是
}

//------------------------------------------------------------------------------
// G2GW_ConnServerProto2  去gateway去链接
type G2GW_ConnServer struct {
	Protocol  int
	Protocol2 int
	ServerID  string //全局配置 唯一的也是
}

// GW2G_ConnServerProto2 返回的数据链接
type GW2G_ConnServer struct {
	Protocol  int
	Protocol2 int
	ServerID  string
}

//------------------------------------------------------------------------------
