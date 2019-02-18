package conf

// csv配置表
var G_GameListNew map[string]*GameListNew // 游戏大厅列表New版本
var G_GameList map[string]*GameList       // 游戏大厅列表
var G_BannerList map[string]*Banner       // 游戏轮播列表
var G_RoomList map[string]interface{}     // 房间列表
var RoomListData map[string]*RoomList     // 房间列表
var RoomListDatabak map[string]*RoomList  // 房间列表
var DSQGameExp map[string]*DSQ_Exp        // 斗兽棋经验列表
var G_ServerList map[string]*ServerList   // 服务器列表

func init() {
	G_GameList = make(map[string]*GameList)
	G_GameListNew = make(map[string]*GameListNew)
	G_BannerList = make(map[string]*Banner)
	G_RoomList = make(map[string]interface{})
	RoomListData = make(map[string]*RoomList)
	RoomListDatabak = make(map[string]*RoomList)
	DSQGameExp = make(map[string]*DSQ_Exp)
	G_ServerList = make(map[string]*ServerList)
	return
}

//------------------------------------------------------------------------------

// 游戏列表
type GameListNew struct {
	GameID    string
	Name      string
	IconPath  string
	IsShelves string
	StartTime string
	EndTime   string
	IsNewest  string
	IsHot     string
	ResPath   string
}

//------------------------------------------------------------------------------
type ServerList struct {
	ID      string
	Name    string
	IP_Port string
}

//------------------------------------------------------------------------------

// 斗兽棋 游戏等级列表
type DSQ_Exp struct {
	Level string
	Exp   string
}

//------------------------------------------------------------------------------

// 房间列表
type RoomList struct {
	RoomID    string // 房间列表
	NeedPiece string // 进场花费的金币
	NeedLev   string // 进场需要的等级
	Desc      string // 描述
	SysPiece  string // 系统抽水
	WinReward string // 每局获得
	IsTop     string // 是否置顶
	TypeICON  string // 活动的ICON
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
// 道具类型
const (
	ITEMTTYPE = iota // ITEMTTYPE == 0
	ItemType1        // ItemType1 == 1 代表货币
	ItemType2        // ItemType1 == 2 代表门票
	ItemType3        // ItemType1 == 3 代表兑换
	ItemType4        // ItemType1 == 4 代表道具
)

// 道具表
type ItemList struct {
	ItemID    string
	ItemName  string
	ItemType  int
	ItemICON  string
	ItemDesc  string
	ItemCoin  string // 兑换的钻石的数量
	IsLimTime string // 是否限时
	LimTime   string // 限时时间
	IsUse     string // 是否可以直接使用
}

//------------------------------------------------------------------------------
// 兑换列表
type AwardList struct {
	AwardID   string // 前端在玩家兑换中查找
	AwardName string
}

//------------------------------------------------------------------------------
