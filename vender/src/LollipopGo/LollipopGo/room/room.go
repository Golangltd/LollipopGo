package room

import (
	"LollipopGo/LollipopGo/player"
)

/*
  房间信息：
	1 房间的匹配机制
	2 房间的总体管理 ，包括：房间的创建、房间的牌型的初始化
	3 房间的数的销毁
	4 房间的机器人的管理等（非重点），机器人完善后就可以增加。

*/

type roominfo interface {
	CreateRoom() *RoomST
	DestroyRoom(string)
}

type RoomST struct {
	RoomID     string
	RoomName   string
	RoomPlayer map[string]*player.PlayerSt // 房间内的玩家
	AllTime    string                      // 时间戳
}

// 房间的管理
// 1 房间的生成，房间的销毁
// 2 人物的匹配
type RoomManger struct {
	RoomInfo map[string]*RoomST
}
