package Mysyl_DB

import (
	"LollipopGo/LollipopGo/player"
	_ "Proto/Proto2"
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
// 设计到30封邮件
// select  * from `t_playeremail` WHERE state = 1 order by itime desc LIMIT 30
func (this *mysql_db) ReadPlayerEmailInfoData(stropenid string) map[int]*player.EmailST {
	rows, err := this.STdb.Query("SELECT * FROM t_playeremail WHERE state = 1 and openid = " + stropenid + " order by itime desc LIMIT 30")
	defer rows.Close()
	CheckErr(err)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("没有错误!")
	}

	dataEmail := make(map[int]*player.EmailST)
	for rows.Next() {
		PlayerSt := new(player.EmailST)
		var uid, times string = "", ""
		rows.Scan(&PlayerSt.ID, &PlayerSt.Name, &PlayerSt.Sender, &PlayerSt.Type,
			&PlayerSt.Time, &PlayerSt.Content, &PlayerSt.IsAdd_ons,
			&PlayerSt.ItemList, &uid,
			&times)

		dataEmail[PlayerSt.ID] = PlayerSt
	}

	fmt.Println("++++++++++++====", dataEmail)

	return dataEmail
}

//------------------------------------------------------------------------------
func (this *mysql_db) ReadAdminEmailInfoData() map[int]*player.EmailST {
	rows, err := this.STdb.Query("SELECT * FROM t_adminemail where state = 1")
	defer rows.Close()
	CheckErr(err)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("没有错误!")
	}

	dataEmail := make(map[int]*player.EmailST)
	for rows.Next() {
		PlayerSt := new(player.EmailST)
		var uid, times string = "", ""
		rows.Scan(&PlayerSt.ID, &PlayerSt.Name, &PlayerSt.Sender, &PlayerSt.Type,
			&PlayerSt.Time, &PlayerSt.Content, &PlayerSt.IsAdd_ons,
			&PlayerSt.ItemList, &uid,
			&times)

		dataEmail[PlayerSt.ID] = PlayerSt
	}

	fmt.Println("++++++++++++====", dataEmail)
	// 更新 标志位 state = 0
	this.Modefy_AdminGameEmailInfoDataGM()
	if len(dataEmail) > 0 {
		// 更新所有的玩家的数据字段 ---> 玩家openid的列表维护一个，先更新内存，再更新数据库
	}
	return dataEmail
}

//------------------------------------------------------------------------------
func (this *mysql_db) ReadUserGameExpInfoData(openid string) int {
	rows, err := this.STdb.Query("SELECT gameexp FROM t_usergameinfo  where openid = " + openid)
	defer rows.Close()
	CheckErr(err)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("没有错误!")
	}
	for rows.Next() {
		var icount int = 0
		rows.Scan(&icount)
		return icount
	}

	return 0
}

//------------------------------------------------------------------------------
// 查询表  select 1 from tablename where uid = 'uid' limit 1;
func (this *mysql_db) ReadUserGameInfoData(openid string) bool {
	rows, err := this.STdb.Query("SELECT * FROM t_usergameinfo  where openid = " + openid + " limit 1")
	defer rows.Close()
	CheckErr(err)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("没有错误!")
	}
	bret := false
	for rows.Next() {
		fmt.Print("查询玩家的游戏列表有数据!!!")
		bret = true
	}

	return bret
}

//------------------------------------------------------------------------------
func (this *mysql_db) ReadUserInfoDataByOpenID(data *player.PlayerSt) (bool, player.PlayerSt) {
	rows, err := this.STdb.Query("SELECT * FROM t_userinfo  where openid = " + data.OpenID + " limit 1")
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
		eer := rows.Scan(&PlayerSt.UID, &uid, &PlayerSt.OpenID, &PlayerSt.VIP_Lev,
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

//------------------------------------------------------------------------------
// 查询表  select 1 from tablename where uid = 'uid' limit 1;
func (this *mysql_db) ReadUserInfoData(uid string) (bool, player.PlayerSt) {
	rows, err := this.STdb.Query("SELECT * FROM t_userinfo  where uid = " + uid + " limit 1")
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
		eer := rows.Scan(&PlayerSt.UID, &uid, &PlayerSt.OpenID, &PlayerSt.VIP_Lev,
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
