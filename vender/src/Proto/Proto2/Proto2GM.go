package Proto2

import (
	"LollipopGo/LollipopGo/player"
)

// G_GameGM_Proto  == 11 子协议
const (
	GAMEINIT                      = iota // GAMEINIT == 0
	W2GMS_Modify_PlayerDataProto2        // W2GMS_Modify_PlayerDataProto2 == 1 修改玩家的数据 :web请求 GM 系统
	GMS2W_Modify_PlayerDataProto2        // GMS2W_Modify_PlayerDataProto2 == 2

	/*
	   邮件 and 跑马灯
	*/
	W2GMS_Modify_PlayerEmailDataProto2 // W2GMS_Modify_PlayerEmailDataProto2  == 3 修改邮件数据
	GMS2W_Modify_PlayerEmailDataProto2 // GMS2W_Modify_PlayerEmailDataProto2  == 4
)

//------------------------------------------------------------------------------
/*
   邮件*跑马灯
*/
// W2GMS_Modify_PlayerEmailDataProto2
type W2GMS_Modify_PlayerEmailData struct {
	Protocol  int
	Protocol2 int
	IMsgtype  int             // 1:表示邮件，2：跑马灯消息，3:针对个人
	OpenID    string          // 玩家唯一ID
	EmailData *player.EmailST // 邮件的消息
	MsgData   *player.MsgST   // 跑马灯的消息
}

// GMS2W_Modify_PlayerEmailDataProto2
type GMS2W_Modify_PlayerEmailData struct {
	Protocol  int
	Protocol2 int
	IMsgtype  int // 1:表示邮件，2：跑马灯消息
	ResultID  int // 结果ID，0：表示成功
}

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
	MODIFY_VIP_LEV        //  MODIFY_VIP_LEV == 5  修改VIP等级
)

//------------------------------------------------------------------------------
// W2GMS_Modify_PlayerDataProto2
// 修改玩家的数据的GM 指令
type W2GMS_Modify_PlayerData struct {
	Protocol  int
	Protocol2 int
	UID       string // 玩家的唯一ID信息
	Itype     string // MODIFYINIT 查看枚举
	ModifyNum string // 正数表示增加，负数标识减少
}

// GMS2W_Modify_PlayerDataProto2
// 返回的操作是否成功
type GMS2W_Modify_PlayerData struct {
	Protocol  int
	Protocol2 int
	Isucc     bool
}

//------------------------------------------------------------------------------
