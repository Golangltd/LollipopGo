package MsgHandleClt

/*
ps: v 2.8.X 版本以前
type Msg_data interface {
	HandleCltProtocol(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{}, Connection *websocket.Conn) interface{}
	HandleCltProtocol2(protocol2 interface{}, ProtocolData map[string]interface{}, Connection *websocket.Conn) interface{}
	PlayerSendMessage(senddata interface{}) int
	CloseEOF(closeEvent interface{}) int
}
*/

// v 2.9.X 以后版本
type Msg_data interface {
	HandleCltProtocol(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{}, Connection interface{}) interface{}
	HandleCltProtocol2(protocol2 interface{}, ProtocolData map[string]interface{}, Connection interface{}) interface{}
	PlayerSendMessage(senddata interface{}) int
	CloseEOF(closeEvent interface{}) int
}

type Msg_dataPB interface {
	HandleCltProtocolPB(ProtocolData []byte, Connection interface{}) interface{}
	HandleCltProtocolPB2(protocol int32, protocol2 int32, ProtocolData []byte, Connection interface{}) interface{}
	PlayerSendMessagePB(senddata []byte) int
	CloseEOFPB(closeEvent interface{}) int
}
