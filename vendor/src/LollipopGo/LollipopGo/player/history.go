package player

/*
历史游戏：
1.	游戏图标
2.	游戏名称
3.	游戏等级
4.	历史最高分数
5.	胜率
*/

type HistoryGame struct {
	UID       int    // UID  存储唯一ID
	GameName  string // 游戏名字
	GameICON  string // 游戏ICON  服务器可以不下发
	GameLev   int    // 游戏的等级
	TopRecord int    // 历史分数最高
	WinRate   string // 胜率 48%
}

//------------------------------------------------------------------------------

/*
历史比赛：
1.	比赛图标
2.	比赛名称
3.	比赛名次
4.	获奖时间
*/

type HistoryRace struct {
	UID      int    // UID  存储唯一ID
	RaceICON string // 比赛的图标
	RaceRank int    // 比赛的名次，最好的名次
	RaceTime string // 比赛的时间 2018.12.24
}

//------------------------------------------------------------------------------
