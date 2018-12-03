package Proto2

const (
	ININGATEWAY             = iota
	C2GWS_PlayerLoginProto2 // C2GWS_PlayerLoginProto2 == 1 登陆协议
	S2GWS_PlayerLoginProto2 // S2GWS_PlayerLoginProto2 == 2 登陆协议
)

// 登陆  客户端--> 服务器
type C2GWS_PlayerLogin struct {
	Protocol  int
	Protocol2 int
	Token     string
}

type S2GWS_PlayerLogin struct {
	Protocol  int
	Protocol2 int
	// 返回给用户-- 用户的信息：UID ，名字，等级等
}
