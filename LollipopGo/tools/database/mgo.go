package database

import (
	"github.com/globalsign/mgo"
	"time"
)

func NewMongoSession(host string) *mgo.Session {
	session, err := mgo.DialWithTimeout(host, 3*time.Second)
	if err != nil {
	}
	session.SetPoolLimit(300)
	return session
}
