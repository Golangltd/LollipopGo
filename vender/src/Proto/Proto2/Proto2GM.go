package Proto2

// G_GameGM_Proto  == 11 子协议
const (
	GAMEINIT                      = iota // GAMEINIT == 0
	W2GMS_Modify_PlayerDataProto2        // W2GMS_Modify_PlayerDataProto2 == 1 修改玩家的数据 :web请求 GM 系统
	GMS2W_Modify_PlayerDataProto2        // GMS2W_Modify_PlayerDataProto2 == 1 修改玩家的数据
)

//------------------------------------------------------------------------------
// 修改玩家的枚举
/*
1 增加钻石数量；减少钻石数量
2 增加金币数量；减少金币数量
3 增加福卡数量；减少福卡数量
4 增加玩家游戏大厅等级；降低玩家游戏大厅等级
5 增加玩家游戏等级；降低玩家游戏等级
*/

const (
	MODIFYINIT     = iota //  MODIFYINIT == 0
	MODIFY_COIN           //  MODIFY_COIN == 1     修改金币
	MODIFY_LEV            //  MODIFY_LEV == 2      修改等级,大厅的等级
	MODIFY_MASONRY        //  MODIFY_MASONRY == 3  修改砖石
	MODIFY_MCARD          //  MODIFY_MCARD == 4    修改福卡
)

//------------------------------------------------------------------------------
// W2GMS_Modify_PlayerDataProto2
// 修改玩家的数据的GM 指令
type W2GMS_Modify_PlayerData struct {
	Protocol  int
	Protocol2 int
	UID       int // 玩家的唯一ID信息
	Itype     int // MODIFYINIT 查看枚举
	ModifyNum int // 正数表示增加，负数标识减少
}

// GMS2W_Modify_PlayerDataProto2
// 返回的操作是否成功
type GMS2W_Modify_PlayerData struct {
	Protocol  int
	Protocol2 int
	Isucc     bool
}

//------------------------------------------------------------------------------
