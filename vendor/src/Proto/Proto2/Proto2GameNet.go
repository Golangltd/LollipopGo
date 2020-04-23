package Proto2

// GameNet_Proto 的子协议
const (
	INIT_NetPROTO2           = iota
	Net_HeartBeatProto2      // Net_HeartBeatProto2 == 1      心跳协议
	Net_Kicking_PlayerProto2 // Net_Kicking_PlayerProto2 == 2 踢人
	Net_RelinkProto2         // Net_RelinkProto2 == 3         断线重新链接

)

//  断线重新链接
//  玩家的结构协议  ---- update
type Net_Relink struct {
	Protocol     int    // 主协议 -- 模块化
	Protocol2    int    // 子协议 -- 模块化的功能
	OpenID       string // 玩家的唯一ID
	StrLoginName string //
	StrLoginPW   string // 加密的数据
	ISucc        bool   // 服务器返回的数据
}

// 踢人
type Net_Kicking_Player struct {
	Protocol  int    // 主协议 -- 模块化
	Protocol2 int    // 子协议 -- 模块化的功能
	OpenID    string // 玩家唯一ID
	ErrorCode int    // 错误码 10001 10002 10003
	StrMsg    string // 原因
}

// 心跳协议
type Net_HeartBeat struct {
	Protocol  int    // 主协议 -- 模块化
	Protocol2 int    // 子协议 -- 模块化的功能
	OpenID    string // 65位 玩家的唯一ID -- server ---> client (多数不需验证OpenID)
}
