package main

import (
	"LollipopGo/LollipopGo/conf"
	"LollipopGo/LollipopGo/log"
	"LollipopGo/LollipopGo/player"
	"LollipopGo/LollipopGo/util"
	_ "LollipopGo/ReadCSV"
	"LollipopGo/db/mysql"
	"Proto"
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

//------------------------------------------------------------------------------
// 保存结算数据
func (t *Arith) SavePlayerDataST2DB(args *Proto2.DB_GameOver, reply *bool) error {
	defer func() {
		if err := recover(); err != nil {
			strerr := fmt.Sprintf("%s", err)
			//发消息给客户端
			ErrorST := Proto2.G_Error_All{
				Protocol:  Proto.G_Error_Proto,      // 主协议
				Protocol2: Proto2.G_Error_All_Proto, // 子协议
				ErrCode:   "80006",
				ErrMsg:    "亲，您发的数据的格式不对！" + strerr,
			}
			fmt.Println("Global server 异常错误", ErrorST)
		}
	}()
	if Mysyl_DB.DB != nil {
		data := Mysyl_DB.DB.InsertPlayerGameInfoST2DB(args)
		*reply = data
	} else {
	}
	return nil
}

// -----------------------------------------------------------------------------
type Arith int

// 登录结构 -- login server
type Args struct {
	A, B int
}

//------------------------------------------------------------------------------
// 获取玩家的数据
// 玩家用户保存
func (t *Arith) GetPlayerST2DB(args *player.PlayerSt, reply *player.PlayerSt) error {
	defer func() {
		if err := recover(); err != nil {
			strerr := fmt.Sprintf("%s", err)
			//发消息给客户端
			ErrorST := Proto2.G_Error_All{
				Protocol:  Proto.G_Error_Proto,      // 主协议
				Protocol2: Proto2.G_Error_All_Proto, // 子协议
				ErrCode:   "80006",
				ErrMsg:    "亲，您发的数据的格式不对！" + strerr,
			}
			fmt.Println("Global server 异常错误", ErrorST)
		}
	}()
	// 1 解析数据 *reply = args.A * args.B
	// roleUID := args.UID
	// 2 保存或者更新数据
	if Mysyl_DB.DB != nil {
		_, data := Mysyl_DB.DB.ReadUserInfoDataByOpenID(args)
		*reply = data
	} else {
	}
	return nil
}

//------------------------------------------------------------------------------
// 修改GM系统
func (t *Arith) ModefyPlayerDataGM(args *Proto2.W2GMS_Modify_PlayerData, reply *Proto2.GMS2W_Modify_PlayerData) error {
	//--------------------------------------------------------------------------
	defer func() {
		if err := recover(); err != nil {
			strerr := fmt.Sprintf("%s", err)
			//发消息给客户端
			ErrorST := Proto2.G_Error_All{
				Protocol:  Proto.G_Error_Proto,      // 主协议
				Protocol2: Proto2.G_Error_All_Proto, // 子协议
				ErrCode:   "80006",
				ErrMsg:    "亲，您发的数据的格式不对！" + strerr,
			}
			fmt.Println("GM server 异常错误：", ErrorST)
		}
	}()
	//--------------------------------------------------------------------------
	uid := util.Str2int_LollipopGo(args.UID)
	itype := util.Str2int_LollipopGo(args.Itype)
	modifynum := util.Str2int_LollipopGo(args.ModifyNum)
	//--------------------------------------------------------------------------
	switch itype {
	case Proto2.MODIFY_COIN, Proto2.MODIFY_LEV, Proto2.MODIFY_MASONRY,
		Proto2.MODIFY_MCARD, Proto2.MODIFY_VIP_LEV:
		bret := Mysyl_DB.DB.Modefy_PlayerDataGM(uid, itype, modifynum)
		*reply = Proto2.GMS2W_Modify_PlayerData{
			Protocol:  Proto.G_GameGM_Proto,
			Protocol2: Proto2.GMS2W_Modify_PlayerDataProto2,
			Isucc:     bret,
		}
	default:
		log.Debug("数据类型不存在!")
	}
	//--------------------------------------------------------------------------
	return nil
}

//------------------------------------------------------------------------------

// 玩家用户保存
func (t *Arith) SavePlayerST2DB(args *player.PlayerSt, reply *player.PlayerSt) error {
	defer func() {
		if err := recover(); err != nil {
			strerr := fmt.Sprintf("%s", err)
			//发消息给客户端
			ErrorST := Proto2.G_Error_All{
				Protocol:  Proto.G_Error_Proto,      // 主协议
				Protocol2: Proto2.G_Error_All_Proto, // 子协议
				ErrCode:   "80006",
				ErrMsg:    "亲，您发的数据的格式不对！" + strerr,
			}
			fmt.Println("Global server 异常错误", ErrorST)
		}
	}()
	// 1 解析数据 *reply = args.A * args.B
	// roleUID := args.UID
	// 2 保存或者更新数据
	if Mysyl_DB.DB != nil {
		_, data := Mysyl_DB.DB.InsertPlayerST2DB(args)
		*reply = data
	} else {
	}
	return nil
}

// 登录的时候，返回的数据
func (t *Arith) Muliply(args *Args, reply *Proto2.GL2C_GameLogin) error {
	// *reply = args.A * args.B
	// 组装数据
	data := &player.GateWayList{
		ServerID:   1001,
		ServerName: "大厅服务器",
		// ServerIPAndPort: "gateway.a.babaliuliu.com:8888", // 测试环境
		ServerIPAndPort: "gateway.b.babaliuliu.com:8888", // 本机  test149.babaliuliu.com
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
