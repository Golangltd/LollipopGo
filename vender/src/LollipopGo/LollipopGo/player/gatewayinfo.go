package player

// 玩家列表
type GateWayList struct {
	ServerID        int
	ServverName     string
	ServerIPAndPort string
	State           string
	OLPlayerNum     int // 服务器目前的玩家数量
	MaxPlayerNum    int
}
