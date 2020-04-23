package Proto3_Data

const (
	INITSERVER           = iota //  INITSERVER == 0
	MJ2GW_ConnInitProto2        //  MJ2GW_ConnInitProto2 == 1 网络初始化
	GW2MJ_ConnInitProto2        //  GW2MJ_ConnInitProto2 == 2
)

type MJ2GW_ConnInit struct {
	Protocol  int
	Protocol2 int
	ServerID  string
}
