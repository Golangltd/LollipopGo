package leaf

import (
	_ "LollipopGo/Proxy_Server/Proto"
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/module"
	"time"
)

const (
	// server conf
	PendingWriteNum = 2000
	MaxMsgLen       = 1 * 1024 * 1024 // 最大长度为1M
	HTTPTimeout     = 5 * time.Second
	LenMsgLen       = 4
	MaxConnNum      = 20000

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000
)

//proto文件序列化/反序列化工具，作为一个全局单例
var MsgProcessor = newGameProcessor()

func NewSkeleton() *module.Skeleton {
	skeleton := &module.Skeleton{
		GoLen:              GoLen,
		TimerDispatcherLen: TimerDispatcherLen,
		AsynCallLen:        AsynCallLen,
		ChanRPCServer:      chanrpc.NewServer(ChanRPCLen),
	}
	skeleton.Init()
	return skeleton
}

func NewGate(wsAddr string, chanRPC *chanrpc.Server) *gate.Gate {
	return &gate.Gate{
		MaxConnNum:      MaxConnNum,
		PendingWriteNum: PendingWriteNum,
		MaxMsgLen:       MaxMsgLen,
		WSAddr:          wsAddr,
		HTTPTimeout:     HTTPTimeout,
		LenMsgLen:       LenMsgLen,
		LittleEndian:    false,
		Processor:       MsgProcessor,
		AgentChanRPC:    chanRPC,
	}
}

func CheckAuth(ag gate.Agent) bool {
	if ag == nil {
		return false
	}
	if ag.UserData() == nil {
		ag.Close()
		return false
	}
	return true
}
