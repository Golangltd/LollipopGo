package WeiXin

import (
	"glog-master"
	"io/ioutil"
	"net/http"
)

// 微信结构
type STWeiXinUserInfo struct {
	OpenID     string
	Nickname   string
	Sex        string
	Province   string
	City       string
	Country    string
	Headimgurl string
	Privilege  string
	Unionid    string
}

// www.GoGoGoEdu.Com
// 接受code
// https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect
func HttpGet(code string) {
	// 所有的配置 微信数据 QQ授权 微博授权数据 全部是配置文件或者数据库配置（通过web端配置）
	appid_redirect_uri := "appid=XXXXXX&redirect_uri=XXXX&response_type="
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?" + appid_redirect_uri + code + "&scope=SCOPE&state=STATE#wechat_redirect"
	resp, err := http.Get(url)
	if err != nil {
		glog.Error("获取用户信息出错：", err.Error())
		return
	}
	_ = resp
	// 数据解析
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		glog.Error("获取用户信息出错：", err1.Error())
		return
	}
	fmt.Println("body", body)
	// 解析成 我们的 微信的结构信息  --获取token
	stb := &STtoken{}
	err2 := json.Unmarshal([]byte(body), &stb)
	if err2 != nil {
		glog.Error("err2---", err2.Error())
		return
	} else {
		// 解析成功
	}
	// 通过token 取玩家数据
	resp1, err1 := http.Get("https://api.weixin.qq.com/sns/userinfo?access_token=" + stb.Access_token + "&openid=" + stb.Openid + "&lang=zh_CN")
	if err1 != nil {
		glog.Error("playerdata", err1.Error())
		return
	}
	// 数据的请求
	body2, err2 := ioutil.ReadAll(resp1.Body)
	if err2 != nil {
		glog.Error("playerdata", err2.Error())
		return
	}
	// 真正获取到了玩家的数据了、
	// 按照一个服务器的解析的玩家的结构 进行组装
	stbtmp := &STWeiXinUserInfo{}
	err2 := json.Unmarshal([]byte(body), &stbtmp)
	if err2 != nil {
		glog.Error("err2---", err2.Error())
		return
	} else {
		// 解析成功
		// 数据的 保存 -- 去数据库 redis
		// --- 自己实现== 头像的信息
	}
	return
}
