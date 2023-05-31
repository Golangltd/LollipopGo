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
	RPC       = "rpc"
	TCP       = "tcp"
	KCP       = "kcp"
)

// 连接服务器类型
const (
	CONNINIT   = iota
	ConnProxy  // ConnProxy == 1 主动连接 Proxy服务器，Proxy作为全球服或者区域服
	StartProxy // StartProxy == 2  Proxy服务器使用
)

// 初始化网络
func InitNet(netty string, netdata ...interface{}) interface{} {
	switch netty {
	case WebSocket:
		InitConnectionPB(netdata[0].(*websocket.Conn))
		return IMsgPB
	case RPC:
		InitConnectionRPC(netty) // rpc 不需要返回
	case KCP:
		InitConnectionKCP(netdata[0].(*kcp.UDPSession), netdata[1].(*kcp.Listener))
		return IMsg
	case TCP:
		InitConnectionTCP(netdata[0].(*net.Conn), netdata[1].(*net.Listener))
		return IMsg
	default:
		glog.Info("InitNet is failed,net type is not exist!")
	}
	return nil
}

// 启动网络监听
func Start(url string, route string, conntype int, netty string) {
	switch netty {
	case WebSocket:

	case RPC:
	case KCP:
	case TCP:
	default:
		glog.Info("InitNet is failed,net type is not exist!")
	}
}
