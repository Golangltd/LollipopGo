package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

// DB的数据的信息
var (
	service = "127.0.0.1:1234"
)

func init() {

	// 注册结构体 + 方法 -->
	// 将结构体的方法注册到rpc中
	arith := new(Arith)
	rpc.Register(arith)

	return
}

// 监听操作
func MainListener(strport string) {

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
		jsonrpc.ServeConn(conn)
	}
}

// -------------------------------------------------
type Args struct {
	A, B int
}

func checkError(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, "Usage: %s", err.Error())
	}
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Muliply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A * args.B
	quo.Rem = args.A / args.B
	return nil
}
