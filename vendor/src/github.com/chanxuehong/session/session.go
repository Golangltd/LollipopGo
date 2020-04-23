// session implements a simple memory-based session container.
// @link        https://github.com/chanxuehong/session for the canonical source repository
// @license     https://github.com/chanxuehong/session/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

// session implements a simple memory-based session container.
// version: 1.1.0
//
//  NOTE: Suggestion is the number of cached elements should not exceed 100,000,
//  because a large number of elements to runtime.GC() is a burden.
//  More than 100,000 can consider memcache, redis ...
//
package session

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNotFound  = errors.New("item not found")
	ErrNotStored = errors.New("item not stored")
)

const (
	// Set the below maximum limit is to prevent int64 overflow, for
	// payload.Expiration = time.Now().Unix() + Storage.maxAge.
	//
	// This can work under normal circumstances,
	// unless deliberately modified the system time to
	// 292277026536-12-05 02:18:07 +0000 UTC
	maxAgeLimit int64 = 60 * 365.2425 * 24 * 60 * 60 // seconds, 60 years
)

type payload struct {
	Key        string
	Value      interface{}
	Expiration int64 // unixtime
}

//                           front                          back
//                          +-------+       +-------+       +-------+
// lruList:list.List        |       |------>|       |------>|       |
//                          |payload|<------|payload|<------|payload|
//                          +-------+       +-------+       +-------+
//                              ^               ^               ^
// cache:                       |               |               |
// map[string]*list.Element     |               |               |
//      +------+------+         |               |               |
//      | key  |value-+---------+               |               |
//      +------+------+                         |               |
//      | key  |value-+-------------------------+               |
//      +------+------+                                         |
//      | key  |value-+-----------------------------------------+
//      +------+------+
//
// Principle:
//   1. len(cache) == lruList.Len();
//   2. for Element of lruList, we get cache[Element.Value.(*payload).Key] == Element;
//   3. in the list lruList, the younger element is always
//      in front of the older elements;
//

// NOTE: Storage is safe for concurrent use by multiple goroutines.
type Storage struct {
	maxAge              int64 // Element of lruList effective time
	gcIntervalResetChan chan time.Duration

	mutex   sync.Mutex
	lruList *list.List
	cache   map[string]*list.Element
}

// New returns an initialized Storage.
//  maxAge:     seconds, the max age of item in Storage; the maximum is 60 years, and can not be adjusted.
//  gcInterval: seconds, GC interval; the minimal is 1 second.
func New(maxAge, gcInterval int) (storage *Storage) {
	_maxAge := int64(maxAge)
	switch {
	case _maxAge <= 0:
		panic(fmt.Sprintf("maxAge must be > 0 and now == %d", _maxAge))
	case _maxAge > maxAgeLimit:
		panic(fmt.Sprintf("maxAge must be <= %d and now == %d", maxAgeLimit, _maxAge))
	}

	if gcInterval <= 0 {
		panic(fmt.Sprintf("gcInterval must be > 0 and now == %d", gcInterval))
	}

	storage = &Storage{
		maxAge:              _maxAge,
		gcIntervalResetChan: make(chan time.Duration),
		lruList:             list.New(),
		cache:               make(map[string]*list.Element, 64),
	}

	// new goroutine for gc service
	go func() {
		gcIntervalDuration := time.Duration(gcInterval) * time.Second

	NEW_GC_INTERVAL:
		ticker := time.NewTicker(gcIntervalDuration)
		for {
			select {
			case gcIntervalDuration = <-storage.gcIntervalResetChan:
				ticker.Stop()
				goto NEW_GC_INTERVAL
			case <-ticker.C:
				storage.gc()
			}
		}
	}()

	return
}

// Get the item cached in the Storage number.
func (storage *Storage) Len() (n int) {
	storage.mutex.Lock()
	n = storage.lruList.Len()
	storage.mutex.Unlock()
	return
}

// Set the NEW gc interval of Storage, seconds.
//  NOTE: if gcInterval <= 0, we do nothing;
//  after calling this method we will soon do a gc(), please avoid business peak.
func (storage *Storage) SetGCInterval(gcInterval int) {
	if gcInterval > 0 {
		storage.gcIntervalResetChan <- time.Duration(gcInterval) * time.Second
		storage.gc() // call gc() immediately
	}
}

// add key-value to Storage.
// ensure that there is no the same key in Storage
func (storage *Storage) add(key string, value interface{}, timeNow int64) (err error) {
	// Check the last element of lruList expired; if so, reused it
	if e := storage.lruList.Back(); e != nil {
		if payload := e.Value.(*payload); timeNow > payload.Expiration {
			delete(storage.cache, payload.Key)

			payload.Key = key
			payload.Value = value
			payload.Expiration = timeNow + storage.maxAge

			storage.cache[key] = e
			storage.lruList.MoveToFront(e)
			return
		}
	}

	// Check whether storage.lruList.Len() has reached the maximum of int.
	if storage.lruList.Len()<<1 == -2 {
		return ErrNotStored
	}

	// now create new
	storage.cache[key] = storage.lruList.PushFront(&payload{
		Key:        key,
		Value:      value,
		Expiration: timeNow + storage.maxAge,
	})
	return
}

// remove Element e from Storage.lruList.
// ensure that e != nil and e is an element of list lruList.
func (storage *Storage) remove(e *list.Element) {
	delete(storage.cache, e.Value.(*payload).Key)
	storage.lruList.Remove(e)
}

// Add key-value to Storage.
// if there already exists a item with the same key, it returns ErrNotStored.
func (storage *Storage) Add(key string, value interface{}) (err error) {
	timeNow := time.Now().Unix()

	storage.mutex.Lock()
	if e, hit := storage.cache[key]; hit {
		if payload := e.Value.(*payload); timeNow > payload.Expiration {
			// payload.Key = key
			payload.Value = value
			payload.Expiration = timeNow + storage.maxAge
			storage.lruList.MoveToFront(e)

			storage.mutex.Unlock()
			return

		} else {
			err = ErrNotStored

			storage.mutex.Unlock()
			return
		}

	} else {
		err = storage.add(key, value, timeNow)

		storage.mutex.Unlock()
		return
	}
}

// Set key-value, unconditional
func (storage *Storage) Set(key string, value interface{}) (err error) {
	timeNow := time.Now().Unix()

	storage.mutex.Lock()
	if e, hit := storage.cache[key]; hit {
		payload := e.Value.(*payload)

		// payload.Key = key
		payload.Value = value
		payload.Expiration = timeNow + storage.maxAge
		storage.lruList.MoveToFront(e)

		storage.mutex.Unlock()
		return

	} else {
		err = storage.add(key, value, timeNow)

		storage.mutex.Unlock()
		return
	}
}

// Get the element with key.
// if there is no such element with the key or the element with the key expired
// it returns ErrNotFound.
func (storage *Storage) Get(key string) (value interface{}, err error) {
	timeNow := time.Now().Unix()

	storage.mutex.Lock()
	if e, hit := storage.cache[key]; hit {
		if payload := e.Value.(*payload); timeNow > payload.Expiration {
			storage.remove(e) // NOTE

			err = ErrNotFound

			storage.mutex.Unlock()
			return

		} else {
			payload.Expiration = timeNow + storage.maxAge
			storage.lruList.MoveToFront(e)

			value = payload.Value

			storage.mutex.Unlock()
			return
		}

	} else {
		err = ErrNotFound

		storage.mutex.Unlock()
		return
	}
}

// Delete the element with key.
// if there is no such element with the key or the element with the key expired
// it returns ErrNotFound, normally you can ignore this error.
func (storage *Storage) Delete(key string) (err error) {
	timeNow := time.Now().Unix()

	storage.mutex.Lock()
	if e, hit := storage.cache[key]; hit {
		storage.remove(e)

		if payload := e.Value.(*payload); timeNow > payload.Expiration {
			err = ErrNotFound
		}

		storage.mutex.Unlock()
		return

	} else {
		err = ErrNotFound

		storage.mutex.Unlock()
		return
	}
}

func (storage *Storage) gc() {
	timeNow := time.Now().Unix()

	storage.mutex.Lock()
	for e := storage.lruList.Back(); e != nil; e = storage.lruList.Back() {
		if payload := e.Value.(*payload); timeNow > payload.Expiration {
			storage.remove(e)
		} else {
			break
		}
	}
	storage.mutex.Unlock()
}
