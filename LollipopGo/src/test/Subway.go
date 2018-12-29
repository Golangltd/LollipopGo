package main

/*
  游戏网关服务器：
     1. 游戏网关服务器可以作为客户端与game server的隔离作用
	 2. 同时可以作为“全球”广播所有玩家的功能
	 3. 消息的加、解密等
     4. 网关可以结合redis作为负载均衡使用，实现分流
	 5. 网关支持TCP、websocket、http等，按需进行分配
*/

//
func init() {
	return
}

//------------------------------------------------------------------------------

func main() {
	strport := "8888"
	http.HandleFunc("/GolangLtd", IndexHandlerGM)
	http.ListenAndServe(":"+strport, nil)
	return
}

func IndexHandlerGM(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}
