package match

import (
	"LollipopGo/LollipopGo/log"
)

var PoolMax int                          // 匹配池的大小
var MapPoolMatch chan map[int]*PoolMatch // 申请数据库

type PoolMatch struct {
	PlayerUID   int // 玩家的UID信息
	MatchTime   int // 玩家匹配的耗时  --- conf配置 数据需要
	PlayerScore int // 玩家的分数
}

func newPoolMatch(IMax int) (MapPoolMatch chan map[int]*PoolMatch) {

	if IMax <= 0 {
		IMax = 100
	}
	return make(chan map[int]*PoolMatch, IMax)
}

func (this *PoolMatch) PutMatch(data map[int]*PoolMatch) {
	if len(MapPoolMatch) >= PoolMax {
		log.Debug("超过了 pool的上限!")
		return
	}
	MapPoolMatch <- data
}

func (this *PoolMatch) GetMatchResult() {}

func (this *PoolMatch) GetMatchNum() {}

func (this *PoolMatch) TimerMatch() {}

func (this *PoolMatch) DestroyMatch() {}
