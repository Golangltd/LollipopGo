package Mysyl_DB

import (
	"database/sql"
)

var (
	dbhostsip      = ""
	dbusername     = ""
	dbpassowrd     = ""
	dbname         = ""
	DB             *mysql_db
	dbMaxOpenConns = 2000
	dbMaxIdleConns = 1000
)

type mysql_db struct {
	STdb *sql.DB
}
