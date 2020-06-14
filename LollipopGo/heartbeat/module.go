package heartbeat

import (
	Proto_Proxy "LollipopGo/Proxy_Server/Proto"
	_ "LollipopGo/Proxy_Server/Proto"
	"LollipopGo/leaf"
	"LollipopGo/tools/tz"
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/module"
	"reflect"
)

//各游戏通用心跳/对时设计

var (
	skeleton      = leaf.NewSkeleton()
	Module        = new(publicModule)
	ChanRPC       = skeleton.ChanRPCServer
	GameRPC       *chanrpc.Server
	EventUserPing interface{}
)

func RegisterGameRPC(id interface{}, gameRPC *chanrpc.Server) {
	EventUserPing = id
	GameRPC = gameRPC
}

type publicModule struct {
	*module.Skeleton
}

func (m *publicModule) OnInit() {
	m.Skeleton = skeleton
}

func (m *publicModule) OnDestroy() {

}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&Proto_Proxy.Ping{}, pong)
	leaf.MsgProcessor.SetRouter(&Proto_Proxy.Ping{}, ChanRPC)
}

func pong(args []interface{}) {
	ag := args[1].(gate.Agent)
	ag.WriteMsg(&Proto_Proxy.Pong{
		Timestamp: tz.GetNowTsMs(),
	})
	if GameRPC != nil && EventUserPing != nil {
		//callback game
		GameRPC.Go(EventUserPing, ag)
	}
}
