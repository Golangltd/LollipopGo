package conf

// csv配置表
var G_GameList map[string]*GameList // 游戏大厅列表
var G_BannerList map[string]*Banner // 游戏轮播列表

func init() {
	G_GameList = make(map[string]*GameList)
	G_BannerList = make(map[string]*Banner)
	return
}

//------------------------------------------------------------------------------

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

//------------------------------------------------------------------------------

// 轮播广告列表
type Banner struct {
	ADID    string
	PicURL  string
	IsTop   string
	SkipURL string // 跳转的URL
	ReMark  string // 备注
}

//------------------------------------------------------------------------------
