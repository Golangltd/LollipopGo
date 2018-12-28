package network

/*
  网络模块
	1. 接口函数
	2. 处理问题的方法，谁去实现？实现接口的结构去实现。
	3. 实现多态
*/

type Conner interface {
	ConnGateWayServer(data interface{})
	PlayerSendMessage(data interface{})
	HandleCltProtocol(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{})
	HandleCltProtocol2(protocol2 interface{}, ProtocolData map[string]interface{})
	Close()
	Destroy()
}

// 注册绑定接口
func RegisterBind(conner Conner) {
	return
}
