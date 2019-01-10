package player

/*
  邮件系统
*/

// 邮件类型
const (
	MSGINIT  = iota // MSGINIT ==0
	MsgType1        // MsgType1 == 1 正常
	MsgType2        // MsgType2 == 2 活动
	MsgType3        // MsgType3 == 3 系统
	MsgType4        // MsgType4 == 4 置顶
)

type EmailST struct {
	ID        int
	Sender    string
	Name      string
	Type      int
	Time      int
	Content   string
	IsAdd_ons bool // 是否有附件
	IsOpen    bool // 是否打开过
	IsGet     bool // 是否打开过
	ItemList  map[int]*ItemST
}
