package Proto2

// G_GameRace_Proto  == 13  子协议
const (
	GameRaceINIT                    = iota // 协议初始化
	C2S_GameRaceBaoMingProto2              // C2S_GameRaceBaoMingProto2 ==1            游戏报名
	S2C_GameRaceBaoMingProto2              // S2C_GameRaceBaoMingProto2 ==2            服务器返回数据
	G_Broadcast_GameRaceStartProto2        // G_Broadcast_GameRaceStartProto2 ==3      比赛开始发送比赛的房间信息
	G_Broadcast_GameRaceDataProto2         // G_Broadcast_GameRaceDataProto2 ==4       广播比赛信息，包括正在比赛的玩家和结束的玩家
	// 循环赛：人数限制需要配置,胜利场数排序
	G_Broadcast_GameRaceResultProto2 // G_Broadcast_GameRaceResultProto2 ==5           循环赛结果全部发给每个客户端，按照胜利次数排行。
	// 积分赛：
	// 段位赛：
)

//------------------------------------------------------------------------------
// G_Broadcast_GameRaceResultProto2  循环赛结果
type G_Broadcast_GameRaceResult struct {
	Protocol   int
	Protocol2  int
	RaceResult map[string]interface{} // 比赛结果
}

//------------------------------------------------------------------------------
// G_Broadcast_GameRaceDataProto2   广播比赛信息，包括正在比赛的玩家和结束的玩家
type G_Broadcast_GameRaceData struct {
	Protocol  int
	Protocol2 int
	Racing    int // 比赛中
	Raced     int // 结束
}

//------------------------------------------------------------------------------
//  G_Broadcast_GameRaceStartProto2   比赛开始
type G_Broadcast_GameRaceStart struct {
	Protocol  int
	Protocol2 int
	RoomTmpID int // 自动匹配成功的房间ID信息
	OpenIDA   string
	OpenIDB   string
}

//------------------------------------------------------------------------------
//C2S_GameRaceBaoMingProto2   游戏报名
type C2S_GameRaceBaoMing struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	Itype     int // 1:参加循环赛，2：参加积分赛，3：参加段位赛
}

//S2C_GameRaceBaoMingProto2   游戏报名
type S2C_GameRaceBaoMing struct {
	Protocol  int
	Protocol2 int
	IState    int // 状态码 1：报名成功，2：报名已经结束，3：其他错误
}

//------------------------------------------------------------------------------
