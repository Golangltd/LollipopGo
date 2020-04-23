package Proto3_Data

const (
	INITUSER             = iota //  INITUSER == 0
	C2GW_ConnLoginProto2        //  C2GW_ConnLoginProto2 == 1 玩家登陆协议
	GW2C_ConnLoginProto2        //  GW2C_ConnLoginProto2 == 2

	MJS2GW_ConnLoginProto2 //  MJS2GW_ConnLoginProto2 == 3 游戏服务器登陆协议
	GW2MJS_ConnLoginProto2 //  GW2MJS_ConnLoginProto2 == 4
)

//------------------------------------------------------------------------------

//------------------------------------------------------------------------------
// 用户登陆协议
type GW2C_ConnLogin struct {
}

type C2GW_ConnLogin struct {
}

//------------------------------------------------------------------------------
