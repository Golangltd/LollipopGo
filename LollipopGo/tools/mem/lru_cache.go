package mem

import (
	"github.com/pkg/errors"
	"sync"
)

var (
	KeyError = errors.New("Key not exist")
)
//缓存在内存里的lru cache
type node struct {
	Prev *node
	Next *node
	Key  string
	Val  interface{}
}

func (nd *node) Detach() {
	if nd.Prev != nil {
		nd.Prev.Next = nd.Next
	}
	if nd.Next != nil {
		nd.Next.Prev = nd.Prev
	}
}

type LRU struct {
	head     *node	//头部节点
	tail     *node	//尾部节点
	keyNodes map[string]*node //key统一使用字符串
	capacity int  //缓存大小
	sync.Mutex
}

//maxSize: 缓存的最大数量
func NewLRUCache(capacity int) *LRU {
	if capacity <= 0 {
		capacity = 40000
	}
	result := &LRU{
		head:     &node{},
		tail:     &node{},
		keyNodes: make(map[string]*node, capacity),
		capacity: capacity,
	}
	result.head.Prev = nil
	result.head.Next = result.tail
	result.tail.Prev = result.head
	result.tail.Next = nil
	return result
}

//将节点移动到链表最前面，在外围加锁
func (lru *LRU) mvHead(nd *node) {
	if lru.head.Next == nd {
		return
	}
	nd.Detach()
	nd.Next = lru.head.Next
	nd.Prev = lru.head
	lru.head.Next = nd
	nd.Next.Prev = nd
}

func (lru *LRU) Get(key string) (interface{}, error) {
	lru.Lock()
	defer lru.Unlock()
	if node, ok := lru.keyNodes[key]; ok {
		lru.mvHead(node)
		return node.Val, nil
	} else {
		return nil, KeyError
	}
}

//删除尾部节点，外围加锁
func (lru *LRU) rmTail() {
	oldTail := lru.tail.Prev
	delete(lru.keyNodes, oldTail.Key)
	oldTail.Detach()
}

func (lru *LRU) Set(key string, value interface{}) {
	lru.Lock()
	defer lru.Unlock()
	if nd, ok := lru.keyNodes[key]; ok {
		lru.mvHead(nd)
		nd.Val = value
	} else {
		nd = &node{
			Key: key,
			Val: value,
		}
		if len(lru.keyNodes) == lru.capacity {
			lru.rmTail()
		}
		lru.keyNodes[key] = nd
		lru.mvHead(nd)
	}
}
