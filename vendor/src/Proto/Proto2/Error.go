package Proto2

// Error_Proto 的子协议
const (
	Error_PROTO2 = iota

	G_Error_All_Proto // G_Error_All_Proto == 1    错误

)

//  错误处理
type G_Error_All struct {
	Protocol  int    // 主协议 -- 模块化
	Protocol2 int    // 子协议 -- 模块化的功能
	ErrCode   string // 玩家的结构
	ErrMsg    string
}
