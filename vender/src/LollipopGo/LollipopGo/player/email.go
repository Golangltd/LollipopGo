package player

/*
  邮件系统
*/

// 邮件类型
const (
	MSGINIT1  = iota // MSGINIT ==0
	MsgType11        // MsgType1 == 1 正常
	MsgType21        // MsgType2 == 2 活动
	MsgType31        // MsgType3 == 3 系统
	MsgType41        // MsgType4 == 4 置顶
)

type EmailST struct {
	ID        int
	Name      string
	Sender    string
	Type      int
	Time      int
	Content   string
	IsAdd_ons bool // 是否有附件
	IsOpen    bool // 是否打开过
	IsGet     bool // 是否打开过
	ItemList  map[int]*ItemST
}

//------------------------------------------------------------------------------

// GM 系统的邮件的结构
type EmailGM struct {
	UID        int
	Name       string
	OPType     int    // 操作类型 1 新增 2 编辑  3 删除
	SendType   int    // 1 表示全服玩家  2 指定玩家发送
	PlayerUID  string // 玩家UID
	SendTime   string // 发送时间，符串1：表示立即发送 ，预定发送表示后面定时发送
	Content    string
	ItemList   string
	EmailState int // 1 已发送 2 冻结  3 发送中（预定） 4 发送失败
}
