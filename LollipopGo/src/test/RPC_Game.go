package main

import (
	"fmt"
)

// 网络接口
type Conn interface {
	// ConnGateWayServer(data interface{})
	// PlayerSendMessage(data interface{})
	// HandleCltProtocol(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{})
	// HandleCltProtocol2(protocol2 interface{}, ProtocolData map[string]interface{})
	// Close()
	Destroy()
}

// 每个类型服务器结构信息 实例1
type CST struct {
	data Conn   //  接口形式  ---- 主要是【注册】形式,所有的都可以处理
	UID  string //  处理
}

// 每个类型服务器结构信息 实例1
type CST1 struct {
	UID string
}

// 调用
func init() {
	data := new(CST)
	NewData(data)
	data1 := new(CST1)
	NewData(data1)
	return
}

// 不用的实现形式
func (this *CST) Destroy() {
	fmt.Println("111")
	return
}

func (this *CST1) Destroy() {
	fmt.Println("222")
	return
}

func main() {
	return
}

// 共用接口

func Show(data Conn) {
	data.Destroy()
}

// new 数据
func NewData(data Conn) {
	data.Destroy()
	return
}

// -----------------------------------------------------------------------------
// 1 接口
type LollipopGoConn interface {
	Destroy()
}

// 2 结构体设计  全局的结构信息
type GlobalServer struct {
	conn    LollipopGoConn
	IP_data string
	Port    string
}

var G_data map[string]*GlobalServer

func init() {
	G_data = make(map[string]*GlobalServer)
}

// 注册：
// 注册服务器的IP、端口等到内存
func Register(cc, ss string, dataer LollipopGoConn) {
	tmp := new(GlobalServer)
	tmp.conn = dataer
	tmp.IP_data = cc
	tmp.Port = dd
	G_data[tmp.IP_data] = tmp
	return
}
