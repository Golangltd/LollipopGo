package Proto

// 主协议 == 规则
const (
	INIT_PROTO             = iota //  INIT_PROTO == 0
	GameData_Proto                //  GameData_Proto == 1      游戏的主协议      game server 协议
	GameDataDB_Proto              //  GameDataDB_Proto == 2    游戏的DB的主协议  db server 协议
	GameNet_Proto                 //  GameNet_Proto == 3       游戏的NET主协议
	G_Error_Proto                 //  G_Error_Proto == 4       游戏的错误处理
	G_Snake_Proto                 //  G_Snake_Proto == 5       贪吃蛇游戏
	G_GateWay_Proto               //  G_GateWay_Proto == 6     网关协议
	G_GameHall_Proto              //  G_GameHall_Proto == 7    大厅协议
	G_GameLogin_Proto             //  G_GameLogin_Proto == 8   登录服务器协议
	G_GameGlobal_Proto            //  G_GameGlobal_Proto == 9  负责全局的游戏逻辑
	G_GameDSQ_Proto               //  G_GameDSQ_Proto == 10    斗兽棋的主协议
	G_GameGM_Proto                //  G_GameGM_Proto == 11     游戏GM管理系统
	G_GamePay_Proto               //  G_GamePay_Proto == 12    游戏支付系统
	G_GameRace_Proto              //  G_GameRace_Proto  == 13  游戏比赛系统
	G_GameBG_Proto                //  G_GameBG_Proto  == 14  新版本的游戏主要协议
	G_GameBGServer_Proto          //  G_GameBGServer_Proto  == 15  新版本的游戏服务器间主要协议
	G_GameXZServer_Proto          //  G_GameBGServer_Proto  == 16  血战
	G_ActivityServer_Proto        //  G_ActivityServer_Proto == 17 活动服务器主协议
	G_GameDDZServer_Proto         //  G_GameDDZServer_Proto  == 18 斗地主
	G_GameTDHServer_Proto         //  G_GameTDHServer_Proto  == 19 推倒胡
	G_GameRobot                   //  G_GameRobot  == 20 机器人
	G_Teenpatti                   //  G_Teenpatti == 21 Teenpatti 协议
	G_A_Bahar                     //  G_A_Bahar == 22 A_Bahar 协议
)
