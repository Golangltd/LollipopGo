package mem

import (
	"sync"
	"time"
)

type ttlNode struct {
	Val        interface{}
	LastUpdate int64
}
type TTL struct {
	alive    int64
	capacity int
	data     map[string]*ttlNode
	sync.Mutex
}

//aliveSeconds: 缓存有效期，capacity: 缓存容量，注意合理估计capacity为未过期元素的最大数量
func NewTTLCache(aliveSeconds int64, capacity int) *TTL {
	return &TTL{
		alive:    aliveSeconds,
		capacity: capacity,
		data:     make(map[string]*ttlNode, capacity),
	}
}
//当缓存容量超出指定值时，会主动扫描一遍过期键
//但是这并不能保证缓存容量减少到指定范围之内
func (ttl *TTL) Set(key string, value interface{}) {
	ttl.Lock()
	defer ttl.Unlock()
	now := time.Now().Unix()
	if nd, ok := ttl.data[key]; ok {
		nd.Val = value
		nd.LastUpdate = now
	} else {
		ttl.data[key] = &ttlNode{
			Val:        value,
			LastUpdate: now,
		}
		if len(ttl.data) >= ttl.capacity {
			for k, v := range ttl.data {
				if v.LastUpdate-now > ttl.alive {
					delete(ttl.data, k)
				}
			}
		}
	}
}

func (ttl *TTL) Get(key string) (interface{}, error) {
	ttl.Lock()
	defer ttl.Unlock()
	if nd, ok := ttl.data[key]; ok {
		if time.Now().Unix()-nd.LastUpdate >= ttl.alive {
			delete(ttl.data, key)
			return nil, KeyError
		} else {
			return nd.Val, nil
		}
	}
	return nil, KeyError
}
