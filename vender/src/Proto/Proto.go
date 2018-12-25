package Proto

// 主协议 == 规则
const (
	INIT_PROTO         = iota //  INIT_PROTO == 0
	GameData_Proto            //  GameData_Proto == 1      游戏的主协议      game server 协议
	GameDataDB_Proto          //  GameDataDB_Proto == 2    游戏的DB的主协议  db server 协议
	GameNet_Proto             //  GameNet_Proto == 3       游戏的NET主协议
	G_Error_Proto             //  G_Error_Proto == 4       游戏的错误处理
	G_Snake_Proto             //  G_Snake_Proto == 5       贪吃蛇游戏
	G_GateWay_Proto           //  G_GateWay_Proto == 6     网关协议
	G_GameHall_Proto          //  G_GameHall_Proto == 7    大厅协议
	G_GameLogin_Proto         //  G_GameLogin_Proto == 8   登录服务器协议
	G_GameGlobal_Proto        //  G_GameGlobal_Proto == 9  负责全局的游戏逻辑
	G_GameDSQ_Proto           //  G_GameDSQ_Proto == 10    斗兽棋的主协议
	G_GameGM_Proto            //  G_GameGM_Proto == 11     游戏GM管理系统
)
