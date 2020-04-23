package match

/*
简单说下实现思路吧：
把参与匹配的玩家都丢入匹配池，每个玩家记录两个属性（分数、开始匹配的时间），每秒
遍历匹配池中所有分数段，找出每个分数上等待时间最长的玩家，用他的范围来进行匹配（因为匹配范围会因为等
待时间边长而增加，等待时间最长的的玩家匹配范围最大，如果连他都匹配不够，那同分数段的其他玩家就更匹配
不够了）。如果匹配到了足够的人，那就把这些人从匹配池中移除，匹配成功；如果匹配人到的人数不够并且没有
达到最大匹配时间，则跳过等待下一秒的匹配；如果达到最大匹配时间，还是没匹配到足够的人，则给这个几个人
凑机器人，提交匹配成功。
*/

type MatchMoudle interface {
	GetMatchResult(string, int) []byte
	PutMatch([]byte)
	GetMatchNum(string) int
	TimerMatch()
	DestroyMatch()
	MatchRecord()
}
