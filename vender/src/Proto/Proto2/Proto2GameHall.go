package Proto2

// G_GameHall_Proto 的子协议
const (
	INIT_HALL = iota

	C2S_HallPlayerLoginProto2 // C2S_HallPlayerLoginProto2 == 1    用户登陆协议
	S2C_HallPlayerLoginProto2 // S2C_HallPlayerLoginProto2 == 2

)

// 客户端-->服务器
// C2S_HallPlayerLoginProto2
type C2S_HallPlayerLogin struct {
	Protocol  int    // 主协议 -- 模块化
	Protocol2 int    // 子协议 -- 模块化的功能
	Token     string // 信息 token
}

// 服务器-->N*客户端(广播出去所有的玩家)
type S2C_HallPlayerLogin struct {
	Protocol  int                   // 主协议 -- 模块化
	Protocol2 int                   // 子协议 -- 模块化的功能
	GameList  map[string]*GGameList // 游戏列表
	// GameListNew map[string]*GGameListNew // 游戏列表
}
