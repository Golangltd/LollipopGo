package database

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

func NewRedisPool(host string, pwd string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp",
				host,
				redis.DialPassword(pwd),
				redis.DialDatabase(db),
				//redis.DialConnectTimeout(5*time.Second),
				//redis.DialReadTimeout(3*time.Second),
				//redis.DialWriteTimeout(3*time.Second),
			)
			if err != nil {
				return nil, fmt.Errorf("can't conn to redis: %s", err)
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
