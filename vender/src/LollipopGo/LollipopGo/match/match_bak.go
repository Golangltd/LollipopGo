package match

import (
	"LollpopGo/LollpopGo/player"
)

var (
	Match_Chan chan *player.PlayerSt
	Imax       int = 0
	ChanMax    int = 1000
	MatchSpeed     = time.Millisecond * 500
)

func init() {
	Match_Chan = make(chan *Match_player, ChanMax)
	go Sort_timer()
}

func Putdata(data *Match_player) {
	Match_Chan <- data
	return
}

func GetChanLength() int {
	Imax = len(Match_Chan)
	return Imax
}

func DoingMatch() {
	Imax = len(MatchChan)
	icount := Imax
	Data := make(map[int]*PlayerST)
	for i := 0; i < Imax; i++ {
		if icount == 1 {
			fmt.Println(MatchChan, "等待匹配")
			continue
		}
		if data, ok := <-MatchChan; ok {
			fmt.Println(data)
			Data[i+1] = data
		} else {
			fmt.Println("woring")
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
