package MsgHandleClt

/*type Msg_data interface {
	HandleCltProtocol(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{}, Connection *websocket.Conn) interface{}
	HandleCltProtocol2(protocol2 interface{}, ProtocolData map[string]interface{}, Connection *websocket.Conn) interface{}
	PlayerSendMessage(senddata interface{}) int
	CloseEOF(closeEvent interface{}) int
}*/


type Msg_data interface {
	HandleCltProtocol(protocol interface{}, protocol2 interface{}, ProtocolData map[string]interface{}, Connection interface{}) interface{}
	HandleCltProtocol2(protocol2 interface{}, ProtocolData map[string]interface{}, Connection interface{}) interface{}
	PlayerSendMessage(senddata interface{}) int
	CloseEOF(closeEvent interface{}) int
}