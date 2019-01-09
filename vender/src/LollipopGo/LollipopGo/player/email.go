package player

/*
  邮件系统
*/

type EmailST struct {
	ID        int
	Name      string
	Time      int
	Content   string
	IsAdd_ons bool // 是否有附件
	IsOpen    bool // 是否打开过
	ItemList  map[int]*ItemST
}
