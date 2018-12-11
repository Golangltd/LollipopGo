package main

import (
	"Proto/Proto2"
	"fmt"
	_ "log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

// DB的数据的信息
var (
	service = "127.0.0.1:8890"
)

func init() {

	// 注册结构体 + 方法 -->
	// 将结构体的方法注册到rpc中
	arith := new(Arith)
	rpc.Register(arith)

	return
}

func checkError(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, "Usage: %s", err.Error())
	}
}

// 监听操作
func MainListener(strport string) {
	arith := new(Arith)
	rpc.Register(arith)
	// 获取数据操作
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+strport)
	checkError(err)

	Listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Fprint(os.Stderr, "accept err: %s", err.Error())
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}

// -----------------------------------------------------------------------------
type Args struct {
	A, B int
}

type Arith int

//------------------------------------------------------------------------------

func (t *Arith) Muliply(args *Args, reply *Proto2.GL2C_GameLogin) error {
	//*reply = args.A * args.B
	*reply = Proto2.GL2C_GameLogin{
		Protocol:  1,
		Protocol2: 2,
		Tocken:    "22222",
		PlayerST:  nil,
		GateWayST: nil,
	}

	return nil
}

// -----------------------------------------------------------------------------
