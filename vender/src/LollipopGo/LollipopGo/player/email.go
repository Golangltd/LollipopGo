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
