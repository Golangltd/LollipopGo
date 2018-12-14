package csv

// csv配置表
var G_StCard2InfoBaseST map[string]*Card2InfoBase // 卡牌活动结构

// 卡牌活动结构
type Card2InfoBase struct {
	Card2ID       string // 卡牌的ID
	Card2Msg      string // 卡牌的描述
	Card2GameName string // 卡牌的地点
	Card2GameID   string // 策划看到的类型
	PicPath       string //  图片路径
	Type          string // 卡牌类型
}
