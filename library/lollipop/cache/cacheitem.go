/*
Golang语言社区(www.Golang.Ltd)
作者：cserli
时间：2018年3月2日
*/
package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	sync.RWMutex

	key interface{}

	data interface{}

	lifeSpan time.Duration

	createdOn time.Time

	accessedOn time.Time

	accessCount int64

	aboutToExpire func(key interface{})
}

func CreateCacheItem(key interface{}, lifeSpan time.Duration, data interface{}) CacheItem {
	t := time.Now()
	return CacheItem{
		key:           key,
		lifeSpan:      lifeSpan,
		createdOn:     t,
		accessedOn:    t,
		accessCount:   0,
		aboutToExpire: nil,
		data:          data,
	}
}

func (item *CacheItem) KeepAlive() {
	item.Lock()
	defer item.Unlock()
	item.accessedOn = time.Now()
	item.accessCount++
}

func (item *CacheItem) LifeSpan() time.Duration {
	// immutable
	return item.lifeSpan
}

func (item *CacheItem) AccessedOn() time.Time {
	item.RLock()
	defer item.RUnlock()
	return item.accessedOn
}

func (item *CacheItem) CreatedOn() time.Time {
	// immutable
	return item.createdOn
}

func (item *CacheItem) AccessCount() int64 {
	item.RLock()
	defer item.RUnlock()
	return item.accessCount
}

func (item *CacheItem) Key() interface{} {
	// immutable
	return item.key
}

func (item *CacheItem) Data() interface{} {
	// immutable
	return item.data
}

func (item *CacheItem) SetAboutToExpireCallback(f func(interface{})) {
	item.Lock()
	defer item.Unlock()
	item.aboutToExpire = f
}
