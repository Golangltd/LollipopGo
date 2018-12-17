package player

// 玩家的结构信息
type PlayerSt struct {
	UID       int
	Name      string
	HeadURL   string
	CoinNum   int
	Awardlist []string // 兑换列表，保存对应的玩家的架构数据里
}
