package csv

import (
	"LollipopGo/LollipopGo/util"
)

var M_CSV *util.Map

func init() {
	M_CSV = new(util.Map)
	ReadCsv_ConfigFile_GameInfoST_Fun()
	ReadCsv_ConfigFile_BannerInfoST_Fun()
	ReadCsv_ConfigFile_RoomListST_Fun()
	ReadCsv_ConfigFile_DSQGameInfoST_Fun()
}
