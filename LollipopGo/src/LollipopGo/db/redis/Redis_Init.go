package Redis_DB

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

//  数据操作
var client *redis.Client

func INIT() {
	// 链接
	client = redis.NewClient(&redis.Options{
		Addr:     "db.a.babaliuliu.com:6379",
		Password: "",
		DB:       0,
	})
	// 心跳  -- 一直做的事情
	// go Redis_Ping(client)

	// return

	// 设置  --  测试时间 耗时 --
	t1 := time.Now()
	for i := 0; i < 1000; i++ {
		err := client.Set("key", i+1, 0).Err()
		if err != nil {
			panic(err)
		}
	}

	elapsed := time.Since(t1)
	fmt.Println("redis ========= Time(s) :", elapsed)

	// 获取
	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

	return
}

// 数据ping数据，保证书的可用性
func Redis_Ping(client *redis.Client) {

	for {
		select {
		case <-time.After(time.Second * 5):
			pong, err := client.Ping().Result()
			if err != nil {
				fmt.Println(err)
			}
			_ = pong
		}
	}
}
