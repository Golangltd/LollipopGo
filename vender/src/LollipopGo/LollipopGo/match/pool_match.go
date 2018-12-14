package match

import (
	"LollipopGo/LollipopGo/log"

	"code.google.com/p/go.net/websocket"
)

var PoolMax int // 匹配池的大小
var MapPoolMatch1V1 chan map[string]*PoolMatch

// 匹配池
type PoolMatch struct {
	OpenID      string          // 玩家的UID加密信息
	MatchType   string          // 1V1 2V2 5V5 等
	Connection  *websocket.Conn // global服务器只是和gateway 进行链接的数据,可以忽略
	MatchTime   int             // 玩家匹配的耗时  --- conf配置 数据需要
	PlayerScore int             // 玩家的分数
}

// 经过算法后的匹配结果
// 1 根据配置的算法进行匹配的操作
type RoomST struct {
	RoomID     string
	RoomName   string
	RoomPlayer map[string]*player.PlayerSt // 房间内的玩家
	AllTime    string                      // 时间戳
}

// 申请链接池
func newPoolMatch(IMax int) (MapPoolMatch chan map[int]*PoolMatch) {

	if IMax <= 0 {
		IMax = 100
	}
	return make(chan map[int]*PoolMatch, IMax)
}

// 玩家点机匹配的时候，需要放入连接池中
func (this *PoolMatch) PutMatch(data map[int]*PoolMatch) {
	// 根据不同的匹配机制，保存不同的数据pool
	if len(MapPoolMatch1V1) >= PoolMax {
		log.Debug("超过了 pool的上限!")
		return
	}
	MapPoolMatch <- data
}

// 根据匹配算法进行返回匹配结果
func (this *PoolMatch) GetMatchResult() {

	// 匹配后就可以发送数据给gateway server 给玩家进行分析
	// send_gateway_data(){}
	return
}

// 获取已经匹配的数量；数量需要记录
// 1 匹配的结果也是需要的发送给DB服务器,玩家登录后返回的数据自带匹配数据
// 2 对战记录
func (this *PoolMatch) MatchRecord() {

	return
}

// 发送数据给gateway server
// 1 这里就是并不需要过多处理
func (this *PoolMatch) PlayerSendMessage() {

	return
}

func (this *PoolMatch) TimerMatch() {}

func (this *PoolMatch) DestroyMatch() {}
