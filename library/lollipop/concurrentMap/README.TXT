安全并发map的使用
func init(){
	M = concurrent.NewConcurrentMap()
}

func XXX(){
		onlineUser := &OnlineUser{
		Connection: ws, // 链接的数据信息== 广播的数据的信息，广播给用户的数据；所有的链接的数据的信息
		MapSafe:    M,  // 并发安全的map
	}
}

func (this *OnlineUser) PullFromClientbb() {

	//------------------------------------------------------------------------------
	// 1 获取唯一标识 链接数据
	val, _ := this.MapSafe.Get("server")
	if val != nil {
		//1  加入 chan

	} else {
		//1  加入map
		this.MapSafe.Put("server", this.Connection)
		//2  起线程
		go this.PullFromClientb()

	}
	// 2 如果map里面存在，就不用开启go;否则开启
	// 3 chan 通信
	//------------------------------------------------------------------------------

	return
}

// 清楚链接的计数器
func GamePanDuanTimer() {
	GGameStartTimerPC := time.NewTicker(20 * time.Second)
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			strerr := fmt.Sprintf("%s", err)
			glog.Info("清除超时PC：", strerr)
			GGameStartTimerPC.Stop() // 暂停计时器
		}
	}()

	for {
		select {
		case <-GGameStartTimerPC.C:
			{
				for itr := M.Iterator(); itr.HasNext(); {
					k, v, _ := itr.Next()
					var key = ""
					var keyName = ""
					// 差分key
					strsplit := Strings_Split(k.(string), "|")
					for i := 0; i < len(strsplit); i++ {
						if len(strsplit) > 2 {
							continue
						}
						if i == 0 {
							key = strsplit[i]
						}
						// 获取链接的名字
						if i == len(strsplit)-1 {
							keyName = strsplit[i]
						}
						// 发消息
						if len(k.(string)) < 41 && keyName == "connect" {

							switch v.(interface{}).(type) {
							case *OnlineUser:
								{
									// 获取数据
									val, _ := v.(interface{}).(*OnlineUser).MapSafe.Get(key + "|MapG_HartJiShu") // 心跳计数
									if val == nil {
										continue
									}
									valG_HartJiShutmp, _ := v.(interface{}).(*OnlineUser).MapSafe.Get(key + "|G_HartJiShutmp")
									if valG_HartJiShutmp == nil {
										continue
									}
									if valG_HartJiShutmp.(int) < val.(int) {
										v.(interface{}).(*OnlineUser).MapSafe.Put(key+"|G_HartJiShutmp", val.(int))
										continue
									}
									valG_HartJiShutmp1, _ := v.(interface{}).(*OnlineUser).MapSafe.Get(key + "|G_HartJiShutmp")
									//									glog.Info("非PC", valG_HartJiShutmp1.(int), "2222wwwwww", val.(int))
									if valG_HartJiShutmp1.(int) >= val.(int) {
										v.(interface{}).(*OnlineUser).MapSafe.Put(key+"|G_HartJiShutmp", 0)
										v.(interface{}).(*OnlineUser).MapSafe.Put(key+"|MapG_HartJiShu", 1)
										v.(interface{}).(*OnlineUser).MapSafe.Remove(key + "|connect")
										v.(interface{}).(*OnlineUser).MapSafe.Remove("server")
										v.(interface{}).(*OnlineUser).Set_GameStart1(false, key)
										glog.Info("清除 超时链接：", key)
									}
								}
							default:
								//glog.Info("非PC")
							}
						}
					}
				}
			}
		}
	}
}

		val, _ := this.MapSafe.Get(strPlayerOpenID + "|connect")
		if val == nil {
			glog.Info("手机已经离线")
		} else {
			val.(interface{}).(*OnlineUser).PlayerSendMessage(GameTiChuInfo)
			this.MapSafe.Remove(strPlayerOpenID + "|connect")
			val.(interface{}).(*OnlineUser).Connection.Close()
		}
		
		
	onlineUser := &OnlineUser{
			Connection: this.Connection, // 链接的数据信息
			StrMD5:     StrFirmUser,     // MD5数据
			MapSafe:    this.MapSafe,
		}
		// 赋值操作数据
		this.MapSafe.Put(code+"|connect", onlineUser)