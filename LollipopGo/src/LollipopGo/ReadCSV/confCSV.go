package csv

// csv配置表
var G_GameList map[string]*GameList // 卡牌活动结构

// 游戏列表
type GameList struct {
	GameID        string // 游戏的ID
	GameName      string // 游戏名字
	GameICON      string // 游戏ICON
	IsShow        string // 是否显示
	ShowStartTime string // 开始显示的时间
	ShowEndTime   string // 结束显示的时间
	IsNewShow     string // 是否最新上架
	IsHotGame     string // 是否是热游戏
}
