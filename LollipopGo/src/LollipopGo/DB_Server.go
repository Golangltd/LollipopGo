package main

import (
	"LollipopGo/LollipopGo/conf"
	"LollipopGo/LollipopGo/player"
	_ "LollipopGo/ReadCSV"
	"LollipopGo/mysql"
	"Proto/Proto2"
	"fmt"
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

// 玩家用户保存
func (t *Arith) SavePlayerST2DB(args *player.PlayerSt, reply *int) error {
	// 1 解析数据 *reply = args.A * args.B
	roleUID := args.UID
	fmt.Println("SavePlayerST2DB:", roleUID)
	// 2 保存或者更新数据
	Mysyl_DB.InsertPlayerST2DB(args)
	return nil

}

// -----------------------------------------------------------------------------
type Arith int

// 登录结构 -- login server
type Args struct {
	A, B int
}

//------------------------------------------------------------------------------

func (t *Arith) Muliply(args *Args, reply *Proto2.GL2C_GameLogin) error {
	//*reply = args.A * args.B
	// 组装数据
	data := &player.GateWayList{
		ServerID:        1001,
		ServerName:      "大厅服务器",
		ServerIPAndPort: "gateway.a.babaliuliu.com:8888",
		State:           "空闲",
		OLPlayerNum:     1024,
		MaxPlayerNum:    5000,
	}
	// 返回数据
	*reply = Proto2.GL2C_GameLogin{
		Protocol:   1,
		Protocol2:  2,
		Tocken:     "22222",
		PlayerST:   nil,
		GateWayST:  data,
		GameList:   conf.G_GameList,
		BannerList: conf.G_BannerList,
	}

	return nil
}

// -----------------------------------------------------------------------------
