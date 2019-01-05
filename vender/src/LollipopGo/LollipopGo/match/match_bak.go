package match

import (
	"LollipopGo/LollipopGo/player"
	"LollipopGo/LollipopGo/util"
	"cache2go"
	"fmt"
	"time"
)

//------------------------------------------------------------------------------

var (
	Match_Chan       chan *player.PlayerSt
	MatchData_Chan   chan map[string]*RoomMatch
	Imax             int = 0
	ChanMax          int = 1000
	MatchSpeed           = time.Millisecond * 500
	PlaterMatchSpeed     = time.Second * 1
	MatchData        map[string]*RoomMatch
	QuitMatchData    map[string]string
	cache            *cache2go.CacheTable
)

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
	MatchData = make(map[string]*RoomMatch)
	MatchData_Chan = make(chan map[string]*RoomMatch, ChanMax)
	QuitMatchData = make(map[string]string)
	cache = cache2go.Cache("myCache")
	go Sort_timer()
}

func Putdata(data *player.PlayerSt) {
	Match_Chan <- data
	return
}

func GetChanLength() int {
	Imax = len(Match_Chan)
	return Imax
}

func DoingMatch() {
	Imax = len(Match_Chan)
	icount := Imax
	Data := make(map[string]*player.PlayerSt)
	iicount := 1
	roomid := ""
	for i := 0; i < Imax; i++ {
		if icount == 1 {
			fmt.Println(Match_Chan, "等待匹配")
			continue
		}
		if data, ok := <-Match_Chan; ok {
			fmt.Println(data)
			if GetMatchPlayer(data.OpenID) {
				fmt.Println(data.OpenID, "玩家已经退出！")
				continue
			}
			Data[util.Int2str_LollipopGo(i+1)] = data
			if iicount%2 == 1 {
				roomid = util.Int2str_LollipopGo(int(util.GetNowUnix_LollipopGo()))
				if MatchData[roomid] == nil {
					continue
				}
				MatchData[roomid].PlayerAOpenID = data.OpenID
			}
			MatchData[roomid].RoomUID = roomid
			MatchData[roomid].RoomLimTime = 10
			MatchData[roomid].RoomPlayerMap[util.Int2str_LollipopGo(i+1)] = Data[util.Int2str_LollipopGo(i+1)]
			MatchData[roomid].PlayerBOpenID = data.OpenID
			if iicount%2 == 0 {
				iicount = 0
				MatchData_Chan <- MatchData
			}
			iicount++
		} else {
			fmt.Println("wrong")
			break
		}
		if icount >= 1 {
			icount--
		}
	}
	if len(Data) > 0 {
		fmt.Println("-------", Data)
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
	cache.Add(OpenID+"QuitMatch", 0, "exit")
}

func DelQuitMatchList(OpenID string) {
	cache.Delete(OpenID + "QuitMatch")
}

func GetMatchPlayer(OpenID string) bool {
	ok := false
	_, err1 := cache.Value(OpenID + "QuitMatch")
	if err1 == nil {
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
}
