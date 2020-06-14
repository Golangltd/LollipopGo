package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

type logger struct {
	Stdout *log.Logger
	StdErr *log.Logger
}

func (lg *logger) Print(values ...interface{}) {
	if values[0] == "sql" {
		return
	}
	msg := fmt.Sprintf("[error  ] db error, msg:%v", values[1:])
	lg.Stdout.Output(3, msg)
	lg.StdErr.Output(3, msg)
}

func newLogger() *logger {
	info := log.New(os.Stdout, "", log.LstdFlags)
	err := log.New(os.Stderr, "", log.LstdFlags)
	return &logger{
		Stdout: info,
		StdErr: err,
	}
}

func NewMysqlConn(host string, isDebug bool) *gorm.DB {
	orm, err := gorm.Open("mysql", host)
	if err != nil {
		log.Fatal("can't init db: ", err)
	}
	orm.DB().SetMaxIdleConns(10)
	orm.DB().SetMaxOpenConns(100)
	orm.SingularTable(true)
	if isDebug {
		orm.LogMode(true)
	} else {
		orm.SetLogger(newLogger())
	}
	return orm
}
