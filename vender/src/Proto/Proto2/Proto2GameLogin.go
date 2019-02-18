package Proto2

import (
	"LollipopGo/LollipopGo/conf"
	"LollipopGo/LollipopGo/player"
)

// G_GameLogin_Proto 的子协议
// 属于HTTP 与 DBserver 进行通信
// 获取到登录正确的信息后，token 返回给网关server
const (
	INIT_GameLogin       = iota
	C2GL_GameLoginProto2 // C2GL_GameLoginProto2 == 1      client -->登录请求
	GL2C_GameLoginProto2 // GL2C_GameLoginProto2 == 2      返回数据
)

//------------------------------------------------------------------------------
// C2GL_GameLoginProto2
// 客户端请求协议
type C2GL_GameLogin struct {
	Protocol  int    // 主协议
	Protocol2 int    // 子协议
	LoginName string // 登录名字
	LoginPW   string // 登录密码
	Timestamp int    // 时间戳
}

// GL2C_GameLoginProto2
// 登录服务器返回给客户端协议
type GL2C_GameLogin struct {
	Protocol    int                          // 主协议
	Protocol2   int                          // 子协议
	Tocken      string                       // server 验证加密数据
	PlayerST    *player.PlayerSt             // 玩家的结构
	GateWayST   *player.GateWayList          // 大厅链接地址
	GameList    map[string]*conf.GameList    // 游戏列表
	GameListNew map[string]*conf.GameListNew // 游戏列表New
	BannerList  map[string]*conf.Banner      // 广告列表
}

//------------------------------------------------------------------------------
