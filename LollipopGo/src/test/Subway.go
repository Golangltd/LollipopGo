package main

/*
  游戏网关服务器：
     1. 游戏网关服务器可以作为客户端与game server的隔离作用
	 2. 消息解析
	 3. 与客户端保持连接，作为广播作用
     4. 消息合法性验证
	 5. 转发消息到业务服务，针对不同的客户端消息做分发到相应的服务处理
	 6. 流量限制，消息分流作用
	 7. 版本验证等
	 8. 可扩展性,动态拓展
*/

////
//func init() {
//	return
//}

////------------------------------------------------------------------------------

//func main() {
//	strport := "8888"
//	http.HandleFunc("/GolangLtd", IndexHandlerGM)
//	http.ListenAndServe(":"+strport, nil)
//	return
//}

//func IndexHandlerGM(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintln(w, "Golang语言社区 www.Golang.Ltd")
//}
