package conf

/*
 配置文件功能：
	1 网络配置，包括启动的IP、端口顺序等
	2 读取游戏中的策划表的配置
	3 集群的配置
*/

var (
	LenStackBuf = 4096

	// log
	LogLevel string
	LogPath  string
	LogFlag  int

	// console
	ConsolePort   int
	ConsolePrompt string = "LollipopGo# "
	ProfilePath   string

	// cluster
	ListenAddr      string
	ConnAddrs       []string
	PendingWriteNum int
)
