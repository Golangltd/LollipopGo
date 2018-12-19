package match

import (
	"LollipopGo/LollipopGo/player"
)

/*
  房间匹配功能：
	1 gameserver 功能，主要是匹配的数据。
	2 房间的数据的管理
	3 定时器的使用等
*/

var MapMatch map[string]*RoomMatch

type RoomMatch struct {
	RoomUID       string                      // 房间号
	RoomName      string                      // 房间名字
	RoomNumPlayer uint8                       // 房间人数
	RoomLimTime   uint64                      // 房间的时间限制
	RoomPlayerMap map[string]*player.PlayerSt // 房间玩家的结构信息
}

func newRoomMatch() (MapMatch map[string]*RoomMatch) {

	return make(map[string]*RoomMatch)
}
