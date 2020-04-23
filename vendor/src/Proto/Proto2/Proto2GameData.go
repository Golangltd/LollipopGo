package Proto2

// GameData_Proto 的子协议
const (
	INIT_PROTO2 = iota

	C2S_PlayerLoginProto2 // C2S_PlayerLoginProto2 == 1    用户登陆协议
	S2C_PlayerLoginProto2 // S2C_PlayerLoginProto2 == 2

	C2S_ChooseRoomProto2 // C2S_ChooseRoomProto2 == 3    选择房间
	S2C_ChooseRoomProto2 // S2C_ChooseRoomProto2 == 4

	C2S_PlayerRunProto2 // C2S_PlayerRunProto2 == 5    模拟走路
	S2C_PlayerRunProto2 // S2C_PlayerRunProto2 == 6

)

// 客户端-->服务器
// C2S_PlayerRunProto2
type C2S_PlayerRun struct {
	Protocol  int    // 主协议 -- 模块化
	Protocol2 int    // 子协议 -- 模块化的功能
	OpenID    string // 角色的ID 信息 token --
	StrRunX   string
	StrRunY   string
	StrRunZ   string // Z
}

// 服务器-->N*客户端(广播出去所有的玩家)
type S2C_PlayerRun struct {
	Protocol  int    // 主协议 -- 模块化
	Protocol2 int    // 子协议 -- 模块化的功能
	OpenID    string // 角色的ID 信息 token --
	StrRunX   string
	StrRunY   string
	StrRunZ   string // Z
}

// -----------------------------------------------------------------------------

// 玩家结构
type PlayerSt struct {
	UID        int
	PlayerName string
	OpenID     string
}

// -----------------------------------------------------------------------------
// 客户端-->服务器
// C2S_PlayerLoginProto2
type C2S_PlayerLogin struct {
	Protocol      int    // 主协议 -- 模块化
	Protocol2     int    // 子协议 -- 模块化的功能
	Itype         int    // 1 登陆，2 代表注册
	Code          string // 微信授权CODE
	StrLoginName  string
	StrLoginPW    string // 加密的数据
	StrLoginEmail string // 收取验证码的
}

// 服务器-->客户端
type S2C_PlayerLogin struct {
	Protocol   int       // 主协议 -- 模块化
	Protocol2  int       // 子协议 -- 模块化的功能
	PlayerData *PlayerSt // 玩家的结构
}

// -----------------------------------------------------------------------------
