package player

// 消息类型
const (
	MSGINIT  = iota // MSGINIT ==0
	MsgType1        // MsgType1 == 1 系统消息
	MsgType2        // MsgType2 == 2 比赛消息
	MsgType3        // MsgType3 == 3 兑奖消息
)

// 系统消息结构
type MsgST struct {
	MsgID     int    // 消息ID
	MsgDesc   string // 消息内容
	MsgState  int    // 消息状态 1，表示上架，2：下架
	MsgType   int    // 1 表示全服
	LoopType  int    // 1：永久，2：次数
	LoopCount int    // 循环次数
	LoopTime  int    // 循环时间
}
