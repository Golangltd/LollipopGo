package log

import (
	"github.com/golang/glog"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
)

/*
 1. 日志收集，设计建议是存储mongoDB数据库
 2. 通信方式采用异步RPC发送到日志服
*/

const (
	LevelINIT = iota
	DebugLevel     // Debug
	ReleaseLevel   // Releases
	ErrorLevel     // Error
	FatalLevel     // Fatal  严重的错误
)

type LogService struct {
	ServerID string
	TypeLog int
	ServiceUrl string
	ConnRPC *rpc.Client
}

type LogSt struct {
	Level int
	Data interface{}
}

func NewLogService(served string,serviceurl string) *LogService {
	client, err := jsonrpc.Dial("tcp", serviceurl)
	if err != nil {
		glog.Info("dial error:", err)
		return nil
	}

	return &LogService{
		ServerID: served,
		TypeLog:  LevelINIT,
		ServiceUrl:serviceurl,
		ConnRPC:client,
	}
}

func (this *LogService)RecordLog(data LogSt)  {

	switch data.Level {
	case DebugLevel:
	case ReleaseLevel:
	case ErrorLevel:
	case FatalLevel:
		log.Fatalln(data)  // 严重错误
	default:
	}
	if this !=nil {
		this.sendlogServer(data)
	}
}

/*
    1. 注册结构
     func rpcRegister() {
		_ = rpc.Register(new(LogSt))
	 }
    2. 函数逻辑
	func (r *LogSt) SaveMongoDB(data log.LogSt, reply *bool) error {
        // 保存数据库逻辑
	}
    3. 基础逻辑

	func main()  {
		conf.InitConfig()
		MainListener(conf.GetConfig().Server.WSAddr)
	}

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

	func checkError(err error) {
		if err != nil {
			fmt.Fprint(os.Stderr, "Usage: %s", err.Error())
		}
	}
    注：日志服务器需要注册 LogSt结构
*/
func (this *LogService)sendlogServer(data LogSt){
	if this.ConnRPC == nil{
		log.Fatalln("初始化错误！")  // 严重错误
		return
	}
	args := data
	var reply bool
	divCall := this.ConnRPC.Go("LogSt.SaveMongoDB", args, &reply, nil)
	replyCall := <-divCall.Done
	glog.Info(replyCall.Reply)
	glog.Info("the LogData return is :", reply)
}