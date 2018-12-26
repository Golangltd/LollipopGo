package player

/*
玩家的结构信息
    注：个人中心所有的数据暂时不修改，数据来自于APP
*/
type PlayerSt struct {
	UID             int                  // 游戏服务器 uid
	VIP_Lev         int                  // VIP 等级
	Name            string               // 玩家的名字
	HeadURL         string               // 玩家的头像
	Sex             string               // 玩家的性别
	PlayerSchool    string               // 学校
	Lev             int                  // 玩家等级
	HallExp         int                  // 玩家大厅的经验
	CoinNum         int                  // 玩家的金币
	MasonryNum      int                  // 玩家的砖石
	MCard           int                  // M 兑换卡
	Constellation   string               // 玩家的星座
	HistoryGameList map[int]*HistoryGame // 历史游戏
	HistoryRaceList map[int]*HistoryRace // 历史比赛
	MedalList       string               // 勋章列表，策划配表
}
