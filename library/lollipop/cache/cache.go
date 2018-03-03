/*
Golang语言社区(www.Golang.Ltd)
作者：cserli
时间：2018年3月2日
*/
package cache

import (
	"sync"
)

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

func Cache(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		t = &CacheTable{
			name:  table,
			items: make(map[interface{}]*CacheItem),
		}

		mutex.Lock()
		cache[table] = t
		mutex.Unlock()
	}

	return t
}
