package impl

import (
	"fmt"
	"github.com/golang/glog"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

// RPC数据结构
type RPCSt struct {
	ServiceUrl string
	SendData interface{}  // 发送的数据
	ReplyData interface{}  // 接受的数据
	ConnRPC *rpc.Client
}

func InitConnectionRPC(Addr string) *RPCSt {
	return &RPCSt{
		ServiceUrl:Addr,
		SendData:nil,
		ReplyData:nil,
		ConnRPC:createClientConn(Addr),
	}
}

func createClientConn(Addr string)  *rpc.Client {
	client, err := jsonrpc.Dial("tcp", Addr)
	if err != nil {
		glog.Info("dial error:", err)
		return nil
	}
	return client
}

func (this *RPCSt)GetClientConnRPC() *rpc.Client  {
	if this.ConnRPC != nil{
		return this.ConnRPC
	}else {
		client, err := jsonrpc.Dial("tcp", this.ServiceUrl)
		if err != nil {
			glog.Info("dial error:", err)
			return nil
		}
		return client
	}
}

// 实际操作信息，
func (this *RPCSt)Send_LollipopGoRPC(data RPCSt) interface{} {
	if this.ConnRPC == nil{
		return nil
	}
	args := data
	var reply RPCSt
	divCall := this.ConnRPC.Go("RPCSt.LollipopGoRPC", args, &reply, nil)
	replyCall := <-divCall.Done
	glog.Info(replyCall.Reply)
	glog.Info("the arith.LollipopGoRPC is :", reply)
	return reply
}

//----------------------------------------------------------------------------------------------------------------------
// RPC 服务器调用

func MainListener(strport string) {
	rpcRegister()
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+strport)
	checkError(err)
	Listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		defer func() {
			if err := recover(); err != nil {
				strerr := fmt.Sprintf("%s", err)
				fmt.Println("异常捕获:", strerr)
			}
		}()
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Fprint(os.Stderr, "accept err: %s", err.Error())
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}

func rpcRegister()  {
	_ = rpc.Register(new(RPCSt))
}

func checkError(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, "Usage: %s", err.Error())
	}
}