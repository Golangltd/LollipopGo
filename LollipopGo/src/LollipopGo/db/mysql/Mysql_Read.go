package Mysyl_DB

import (
	"LollipopGo/LollipopGo/player"
	"database/sql"
	"fmt"
)

/*
   插入数据库数据操作
*/

func QueryFromDB(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM userinfo")
	CheckErr(err)
	if err != nil {
		fmt.Println("error:", err)
	} else {
	}
	for rows.Next() {
		var uid string
		var username string
		var departmentname string
		var created string
		var password string
		var autid string
		CheckErr(err)
		err = rows.Scan(&uid, &username, &departmentname, &created, &password, &autid)
		fmt.Println(autid)
		fmt.Println(username)
		fmt.Println(departmentname)
		fmt.Println(created)
		fmt.Println(password)
		fmt.Println(uid)
	}
}

//------------------------------------------------------------------------------
// 查询表  select 1 from tablename where uid = 'uid' limit 1;
func (this *mysql_db) ReadUserInfoData(uid string) (bool, player.PlayerSt) {
	rows, err := this.STdb.Query("SELECT * FROM t_userinfo_copy  where uid = " + uid + " limit 1")
	defer rows.Close()
	CheckErr(err)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("没有错误!")
	}

	// cols, _ := rows.Columns()
	// for i := range cols {
	// 	fmt.Print(cols[i])
	// 	fmt.Print("\t")
	// }
	bret := false
	var PlayerSt player.PlayerSt
	for rows.Next() {
		var uid, times string = "", ""
		eer := rows.Scan(&PlayerSt.UID, &uid, &PlayerSt.VIP_Lev,
			&PlayerSt.Name, &PlayerSt.HeadURL, &PlayerSt.PlayerSchool,
			&PlayerSt.Sex, &PlayerSt.HallExp, &PlayerSt.CoinNum,
			&PlayerSt.MasonryNum, &PlayerSt.MCard,
			&PlayerSt.Constellation, &PlayerSt.MedalList,
			&times)
		fmt.Println("+++++++++8888", PlayerSt, eer)
		bret = true
	}

	return bret, PlayerSt
}
