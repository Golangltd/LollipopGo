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
	MsgID   int    // 消息ID
	MsgType int    // 前端和消息类型去显示到对应的模块
	MsgDesc string // 消息内容
}
