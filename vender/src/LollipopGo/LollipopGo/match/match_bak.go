package match

import (
	"LollipopGo/LollipopGo/player"
	"LollipopGo/LollipopGo/util"
	"fmt"
	"time"
)

//------------------------------------------------------------------------------

var (
	Match_Chan     chan *player.PlayerSt
	MatchData_Chan chan map[string]*RoomMatch
	Imax           int = 0
	ChanMax        int = 1000
	MatchSpeed         = time.Millisecond * 500
	MatchData      map[string]*RoomMatch
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
			// 30s 就剔除
			continue
		}
		if data, ok := <-Match_Chan; ok {
			fmt.Println(data)
			Data[util.Int2str_LollipopGo(i+1)] = data
			// 获取房间ID信息
			if iicount%2 == 1 {
				roomid = util.Int2str_LollipopGo(int(util.GetNowUnix_LollipopGo()))
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
	timer := time.NewTimer(MatchSpeed)
	for {
		select {
		case <-timer.C:
			{
				DoingMatch()
			}
		}
	}
}
