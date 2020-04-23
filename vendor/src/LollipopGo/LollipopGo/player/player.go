package player

/*
玩家的结构信息
    注：个人中心所有的数据暂时不修改，数据来自于APP
*/
type PlayerSt struct {
	UID             int                    // 游戏服务器 uid
	VIP_Lev         int                    // VIP 等级
	Name            string                 // 玩家的名字
	HeadURL         string                 // 玩家的头像
	Sex             string                 // 玩家的性别
	PlayerSchool    string                 // 学校
	Lev             int                    // 玩家等级
	HallExp         int                    // 玩家大厅的经验
	CoinNum         int                    // 玩家的金币
	MasonryNum      int                    // 玩家的砖石
	MCard           int                    // M 兑换卡
	Constellation   string                 // 玩家的星座
	HistoryGameList map[int]*HistoryGame   // 历史游戏
	HistoryRaceList map[int]*HistoryRace   // 历史比赛
	MedalList       string                 // 勋章列表，策划配表
	OpenID          string                 // MD5数据
	GameData        map[int]*PlayerGameLev // 游戏数据
	IsNewEmail      bool                   // 是否有新邮件
}

/*
   游戏的记录，游戏的等级等
*/

type PlayerGameLev struct {
	UID       int
	GameID    int
	GameLev   int
	GameExp   int
	GameScore int
}

/*
   玩家的状态，需要网关经行保存；
*/

const (
	StateInit    = iota // 初始化
	GateWayState        // 网关
	GlobalServer        // 公共服
	DSQServer           // 斗兽棋服
)

// 玩家的网络状态
type PlayerConnState struct {
	OpenID string
	Istate int // 对应上面的状态
}

// 麻将-玩家个人信息
type PlayerItem struct {
	UID          int     //玩家ID
	OpenID       string  //玩家openID
	Nickname     string  //玩家昵称
	Avatar       string  //玩家头像
	Mobile       string  //手机号
	Wechat       string  //微信号
	Wealth       float64 //财富
	Win 		 int 	 // 输赢概率
	IsAuthIDCard bool    //是否实名认证
	IsBindMObile bool    //是否绑定手机
	IsBindWechat bool    //是否绑定微信
	LoginType    int     //登录类型 1游客登录 2手机登录 3微信登录
	Score        int     // 用户分数 比赛分数
}
