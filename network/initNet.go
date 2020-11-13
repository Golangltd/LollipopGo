package impl

import (
	"github.com/golang/glog"
	"github.com/xtaci/kcp-go"
	"golang.org/x/net/websocket"
	"net"
)

// LollipopGo 支持的网络类型
const (
	WebSocket = "websocket"
	RPC = "rpc"
	TCP = "tcp"
	KCP = "kcp"
)

// 初始化网络
func InitNet( netty string ,netdata ...interface{}) interface{} {
	switch netty {
	case WebSocket:
		InitConnection(netdata[0].(*websocket.Conn))
		return IMsg
	case RPC:
		InitConnectionRPC(netty)  // rpc 不需要返回
	case KCP:
		InitConnectionKCP(netdata[0].(*kcp.UDPSession),netdata[1].(*kcp.Listener))
		return IMsg
	case TCP:
		InitConnectionTCP(netdata[0].(*net.Conn),netdata[1].(*net.Listener))
		return IMsg
	default:
		glog.Info("InitNet is failed,net type is not exist!")
	}
	return nil
}




