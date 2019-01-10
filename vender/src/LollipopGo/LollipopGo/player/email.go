package player

/*
  邮件系统
*/

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
