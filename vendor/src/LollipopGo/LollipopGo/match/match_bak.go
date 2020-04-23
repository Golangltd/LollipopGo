package match

import (
	"LollipopGo/LollipopGo/player"
	"LollipopGo/LollipopGo/util"
	"cache2go"
	"fmt"
	"time"

	"code.google.com/p/go.net/websocket"
)

/*
   作者邮箱 : lihaibin@bytedancing.com
   有疑问，发邮箱
*/
//------------------------------------------------------------------------------

var (
	Match_Chan       chan *player.PlayerSt
	MatchData_Chan   chan map[string]interface{}
	Imax             int = 0
	ChanMax          int = 1000
	MatchSpeed           = time.Millisecond * 500
	PlaterMatchSpeed     = time.Second * 1
	MatchData        map[string]interface{}
	QuitMatchData    map[string]string
	cache            *cache2go.CacheTable
	MatchRoomUID     int = 1000
	TimeOutCount     int = 0
	ConnMatch        *websocket.Conn
)

var TempMatch map[string]string

//------------------------------------------------------------------------------

type RoomMatch struct {
	RoomUID       string                      // 房间号
	PlayerAOpenID string                      // A 阵营的OpenID
	PlayerBOpenID string                      // B 阵营的OpenID
	RoomLimTime   uint64                      // 房间的时间限制
	RoomPlayerMap map[string]*player.PlayerSt // 房间玩家的结构信息
}

//------------------------------------------------------------------------------

func init() {
	Match_Chan = make(chan *player.PlayerSt, ChanMax)
	//	MatchData_Chan = make(chan map[string]*RoomMatch, ChanMax)
	MatchData_Chan = make(chan map[string]interface{}, ChanMax)
	QuitMatchData = make(map[string]string)
	cache = cache2go.Cache("myCache")
	go Sort_timer()
}

func Putdata(data *player.PlayerSt) {
	fmt.Println("加入匹配队列")
	Match_Chan <- data
	return
}

func GetChanLength() int {
	Imax = len(Match_Chan)
	return Imax
}

func DoingMatch() {

	TimeOutCount++
	Imax = len(Match_Chan)
	if Imax == 1 {
		fmt.Println(Match_Chan, "等待匹配")
		data := <-Match_Chan
		if len(QuitMatchData[data.OpenID]) > 0 {
			delete(QuitMatchData, data.OpenID)
			return
		} else {
			Match_Chan <- data
			return
		}
	}

	roomid := ""
	icround := Imax / 2

	MatchData = make(map[string]interface{})
	datamatch := new(RoomMatch)
	datamatch.RoomPlayerMap = make(map[string]*player.PlayerSt)

	for i := 1; i < icround*2+1; i++ {

		if data, ok := <-Match_Chan; ok {
			fmt.Println("3333333333333333333333", data)
			datamatch.RoomLimTime = 10
			roomid = util.Int2str_LollipopGo(MatchRoomUID)
			datamatch.RoomUID = roomid
			datamatch.RoomPlayerMap[data.OpenID] = data
			DelMatchQueue(data.OpenID)
		}

		if i%2 == 0 {
			MatchData[roomid] = datamatch
			MatchData_Chan <- MatchData
			fmt.Println("0------------", MatchData_Chan)
			MatchRoomUID++
		}
	}
}

func Sort_timer() {
	for {
		select {
		case <-time.After(MatchSpeed):
			{
				DoingMatch()
			}
		}
	}
}

func SetQuitMatch(OpenID string) {
	QuitMatchData[OpenID] = OpenID
}

func DelQuitMatchList(OpenID string) {
	delete(QuitMatchData, OpenID)
}

func GetMatchPlayer(OpenID string) bool {
	ok := false
	res, err1 := cache.Value(OpenID + "QuitMatch")
	if err1 != nil {
		return ok
	}
	if res.Data().(string) == "exit" {
		ok = true
	}
	return ok
}

func GetMatchQueue(OpenID string) bool {
	ok := false
	_, err1 := cache.Value(OpenID + "MatchQueue")
	if err1 == nil {
		ok = true
	}
	return ok
}

func SetMatchQueue(OpenID string) {
	cache.Add(OpenID+"MatchQueue", 0, "exit")
	DelQuitMatchList(OpenID)
}

func DelMatchQueue(OpenID string) {
	cache.Delete(OpenID + "MatchQueue")
	QuitMatchData[OpenID] = OpenID
}
