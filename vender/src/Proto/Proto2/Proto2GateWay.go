package Proto2

const (
	ININGATEWAY             = iota
	C2GWS_PlayerLoginProto2 // C2GWS_PlayerLoginProto2 == 1 登陆协议
	S2GWS_PlayerLoginProto2 // S2GWS_PlayerLoginProto2 == 2 登陆协议

	GateWay_HeartBeatProto2 // GateWay_HeartBeatProto2 == 3 心跳协议
)

//------------------------------------------------------------------------------
// 断线重连

//------------------------------------------------------------------------------

// C2GWS_PlayerLoginProto2
// 登陆  客户端--> 服务器
type C2GWS_PlayerLogin struct {
	Protocol  int
	Protocol2 int
	PlayerUID string
	Token     string
}

// S2GWS_PlayerLoginProto2
type S2GWS_PlayerLogin struct {
	Protocol   int
	Protocol2  int
	OpenID     string
	PlayerST   *player.PlayerSt          // 玩家的结构
	GateWayST  *player.GateWayList       // 大厅链接地址
	GameList   map[string]*conf.GameList // 游戏列表
	BannerList map[string]*conf.Banner   // 广告列表
}

//------------------------------------------------------------------------------

// GateWay_HeartBeatProto2
// 心跳协议
type GateWay_HeartBeat struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 65位 玩家的唯一ID -- server ---> client (多数不需验证OpenID)
}
