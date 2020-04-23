package WeiXin

var G_StWeiXinDatatmp map[string]*StWeiXinUserInfo

// 微信结构
type StWeiXinUserInfo struct {
	OpenID        string
	Name          string
	Sex           uint32
	Language      string
	City          string
	Province      string
	Country       string
	HeadUrl       string
	Privilege     string
	IdentifykeySJ string
	IdentifykeyHD string
	Masonry       string //
}

func init() {
	G_StWeiXinDatatmp = make(map[string]*StWeiXinUserInfo)
}
