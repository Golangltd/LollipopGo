package game

/*

 通用性质的实现
 1. new() 可以用


*/

type CGame struct {
	ServerID   int
	ServerName string
	ServerNet  string
}

func NewGame(serverdata CGame) *CGame {
	return &CGame{
		ServerID: serverdata.ServerID,
		ServerID: serverdata.ServerName,
		ServerID: serverdata.ServerNet,
	}
}

func (this *CGame) CreateNewGame() {
	// 1 启动链接 gateway 服务器
	// 2 建立new游戏的监听端口
	// 3 建立new的网络处理、消息分发模块
	// 4 建立定时器模块（完整的游戏链接处理）
	// 5 游戏退出后的资源销毁

	return
}

func (this *CGame) ConnetGatewayServer() {
	// 1 链接成功，触发消息分发处理机制
	// 2 心跳触发
	return
}
