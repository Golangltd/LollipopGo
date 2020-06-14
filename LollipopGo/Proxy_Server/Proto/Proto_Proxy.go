package Proto_Proxy

// 代理协议
// 主协议 1
const (
	INIYPROXY             = iota //  ==0
	C2Proxy_SendDataProto        //  C2Proxy_SendDataProto == 1  客户端发送协议
	Proxy2C_SendDataProto        //  Proxy2C_SendDataProto == 2
	G2Proxy_ConnDataProto        //  G2Proxy_ConnDataProto == 3  服务器链接协议
	Proxy2G_ConnDataProto        //  Proxy2G_ConnDataProto == 4
	G2Proxy_SendDataProto        //  G2Proxy_SendDataProto == 5  服务器发送代理
	Proxy2G_SendDataProto        //  Proxy2G_SendDataProto == 6
	C2Proxy_ConnDataProto        //  C2Proxy_ConnDataProto == 7  客户端连接协议
	Proxy2C_ConnDataProto        //  Proxy2C_ConnDataProto == 8
)

//------------------------------------------------------------------------------
// C2Proxy_ConnDataProto  客户端连接协议
type C2Proxy_ConnData struct {
	Protocol  int
	Protocol2 int
	OpenID    string // 客户端第一次发空
}

// Proxy2C_ConnDataProto
type Proxy2C_ConnData struct {
	Protocol  int
	Protocol2 int
	OpenID    string
}

//------------------------------------------------------------------------------
// G2Proxy_SendDataProto  服务器发送代理
type G2Proxy_SendData struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	Data      interface{}
}

// Proxy2G_SendDataProto
type Proxy2G_SendData struct {
	Protocol  int
	Protocol2 int
	OpenID    string
	Data      interface{}
}

//------------------------------------------------------------------------------
// G2Proxy_ConnDataProto   服务器链接协议
type G2Proxy_ConnData struct {
	Protocol  int
	Protocol2 int
	ServerID  string
}

// Proxy2G_ConnDataProto
type Proxy2G_ConnData struct {
	Protocol  int
	Protocol2 int
}

//------------------------------------------------------------------------------
// C2Proxy_SendDataProto  客户端发送协议
type C2Proxy_SendData struct {
	Protocol  int
	Protocol2 int
	ServerID  string
	Data      interface{} //
}

//type Proxy2GS_InitHall struct {
//	Protocol  int
//	Protocol2 int
//	Token     string
//	OpenID    string
//	LoginType string
//	//IMEI      string
//}

// Proxy2C_SendDataProto
type Proxy2C_SendData struct {
	Protocol  int
	Protocol2 int
	ServerID  string
	Data      interface{}
}

//------------------------------------------------------------------------------
