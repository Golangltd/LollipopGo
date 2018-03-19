初始化

cache = cache2go.Cache("myCache")

	// 写入cache
	cache.Add(strtmpopenid, 0, "exit")
	
	// 删除某一个玩家
func Del_GameOneDataClear(strOpenID string) bool {

	cache.Delete(strOpenID)
	//	if G_MapOpenIdtmp[strOpenID] != nil {
	//		delete(G_MapOpenIdtmp, strOpenID)
	//		return true
	//	}
	return false
}

// 检测客户端在不
func (this *OnlineUser) Check_Client_Is_Connect(strcode string) bool {

	//glog.Info(" Entry Check_Client_Is_Connect")
	if len(strcode) == 0 {
		glog.Error(" Check_Client_Is_Connect ,strcode is nil!!!")
		return false
	}

	ok := false
	res, err1 := cache.Value(strcode)
	if err1 == nil {
		ok = true
		//glog.Info("Found value in cache:", res.Data())
	} else if res == nil {
		glog.Info("Not Found value in cache opneid:", strcode)
	}
	return ok
}
