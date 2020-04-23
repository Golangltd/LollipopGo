package match

import (
	"LollipopGo/LollipopGo/log"

	"LollipopGo/LollipopGo/player"

	"code.google.com/p/go.net/websocket"
)

const (
	INITMATCH = iota // INITMATCH == 0
	Match_1V1        // Match_1V1 == 1
	Match_2V2        // Match_2V2 == 2
	Match_3V3        // Match_3V3 == 3
)

var PoolMax int                             // 匹配池的大小
var MapPoolMatch1V1 chan map[int]*PoolMatch // key 是游戏ID

// 匹配池
type PoolMatch struct {
	OpenID      string          // 玩家的UID加密信息
	MatchType   string          // 1V1 2V2 5V5 等
	Connection  *websocket.Conn // global服务器只是和gateway 进行链接的数据,可以忽略
	MatchTime   int             // 玩家匹配的耗时  --- conf配置 数据需要
	PlayerScore int             // 玩家的分数
	PlayerLev   int             // 玩家等级
	MatchGame   int             // 玩家匹配的游戏
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
// map[int]*PoolMatch int就是游戏ID
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
	MapPoolMatch1V1 <- data
}

// 根据匹配算法进行返回匹配结果
// 条件：那款游戏匹配，匹配类型（1V1）
// 定时器：每个秒就要匹配一次所有数据
func (this *PoolMatch) GetMatchResult(igameid int, imatchtype int) {

	if imatchtype == Match_1V1 {
		// 1V1 匹配
		// 找到游戏
		// data <- MapPoolMatch1V1
		// 排序 lev
		// 生成数据，roomID等
		// 发送给网关
	} else if imatchtype == Match_2V2 {
		// 2V2 匹配
	} else if imatchtype == Match_3V3 {
		// 3V3 匹配
	}
	// 匹配数据发给 GateWay Server
	// send_gateway_data(){}
	return
}

// 获取已经匹配的数量；数量需要记录
// 1 匹配的结果也是需要的发送给DB服务器,玩家登录后返回的数据自带匹配数据
// 2 对战记录
func (this *PoolMatch) MatchRecord() {}

// 发送数据给gateway server
// 1 这里就是并不需要过多处理
func (this *PoolMatch) PlayerSendMessage() {}

func (this *PoolMatch) TimerMatch() {}

func (this *PoolMatch) DestroyMatch() {}
