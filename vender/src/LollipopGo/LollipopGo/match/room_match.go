package match

import (
	"LollipopGo/LollipopGo/player"
)

var MapMatch map[string]*RoomMatch

type RoomMatch struct {
	RoomUID       string                      // 房间号
	RoomName      string                      // 房间名字
	RoomNumPlayer uint8                       // 房间人数
	RoomLimTime   uint64                      // 房间的时间限制
	RoomPlayerMap map[string]*player.PlayerSt // 房间玩家的结构信息
}

func newRoomMatch() (MapMatch map[string]*RoomMatch) {

	return make(MapMatch map[string]*RoomMatch)
}