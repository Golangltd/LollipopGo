package Error

const (
	Coin_lack  = 60000 // 金币不足
	Lev_lack   = 60001 // 等级不够
	IsMatch    = 60002 // 等级不够
	LoginError = 10001 // 登陆失败
)

type ErrorMsg struct {
	Protocol  int
	Protocol2 int
	ErrorID   int
	Errormsg  string
}
