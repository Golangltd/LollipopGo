package impl

import (
	"github.com/xtaci/kcp-go"
	"golang.org/x/net/websocket"
)

// LollipopGo 支持的网络类型
const (
	WebSocket = "websocket"
	RPC = "rpc"
	TCP = "tcp"
	UDP = "udp"
	KCP = "kcp"
	NCN = "ncn"
)


// 初始化网络
func InitNet( netty string ,netdata ...interface{}) interface{} {
	switch netty {
	case WebSocket:
		for _, arg := range netdata {
			InitConnection(arg.(*websocket.Conn))
		}
		return IMsg
	case RPC:
	case KCP:
		InitConnectionKCP(netdata[0].(*kcp.UDPSession),netdata[1].(*kcp.Listener))
		return IMsg
	case TCP:
	case UDP:
	case NCN:
	default:
	}
	return nil
}