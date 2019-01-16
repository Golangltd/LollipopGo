package network

var (
	LoginServerAddr  string
	GateWayAddr      string
	DBServerAddr     string
	GlobalServerAddr string
	GMServerAddr     string
	DSQServerAddr    string // 子游戏服务器配置
)

type Conner interface {
	ConnGateWayServer(data interface{})
	PlayerSendMessage(data interface{})
	HandleCltProtocol(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{})
	HandleCltProtocol2(protocol2 interface{}, ProtocolData map[string]interface{})
	Close()
	Destroy()
}
