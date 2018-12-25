package Proto2

const (
	GAMEINIT                      = iota // GAMEINIT == 0
	W2GMS_Modify_PlayerDataProto2        // W2GMS_Modify_PlayerDataProto2 == 1 修改玩家的数据 :web请求 GM 系统
	GMS2W_Modify_PlayerDataProto2        // GMS2W_Modify_PlayerDataProto2 == 1 修改玩家的数据
)

//------------------------------------------------------------------------------
// 修改玩家的枚举
const (
	MODIFYINIT  = iota //  MODIFYINIT == 0
	MODIFY_COIN        //  MODIFY_COIN == 1 修改金币
	MODIFY_LEV         //  MODIFY_LEV == 2  修改等级
)

//------------------------------------------------------------------------------
// W2GMS_Modify_PlayerDataProto2
// 修改玩家的数据的GM 指令
type W2GMS_Modify_PlayerData struct {
	Protocol  int
	Protocol2 int
	Itype     int // MODIFYINIT 查看枚举
	ModifyNum int
}

// GMS2W_Modify_PlayerDataProto2
// 返回的操作是否成功
type GMS2W_Modify_PlayerData struct {
	Protocol  int
	Protocol2 int
	Isucc     bool
}

//------------------------------------------------------------------------------
