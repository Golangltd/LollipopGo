package room

import (
	"LollipopGo/LollipopGo/match"
	"LollipopGo/LollipopGo/player"
)

/*
  房间信息：
	0 流程 --> client 到达  网关gateway  --> global server 进行玩家数据匹配 --> 网关gateway --> client  --> 网关gateway --> game server
	1 房间的匹配机制
	2 房间的总体管理 ，包括：房间的创建、房间的牌型的初始化
	3 房间的数的销毁
	4 房间的机器人的管理等（非重点），机器人完善后就可以增加。
*/

type roominfo interface {
	CreateRoom() *RoomST
	DestroyRoom(string)
}

// 房间的管理
// 1 房间的生成，房间的销毁
// 2 人物的匹配
type RoomManger struct {
	RoomInfo  map[string]*match.RoomST // 已经匹配的玩家的结构
	AllPlayer match.MapPoolMatch       // 匹配的池子
}

// 创建房间
func (this *RoomManger) CreateRoom() {
	// 玩家匹配的机制

	return
}
