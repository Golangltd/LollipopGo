package race

var (
	IRace_Time int
)

// 比赛接口
type ReacIF interface {
	BaoMingRaceData()
	BaoMingExit()
	GetRaceDataFromDB()
	PutRaceDataToDB()
}

// 比赛结构
type RaceSt struct {
	RaceUID  int
	RaceName string
	RaceTime string
	RaceDesc string
	RaceNum  int
	RaceJl   string
}

// 初始化
func init() {
	return
}
