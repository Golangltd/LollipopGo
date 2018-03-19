package concurrent

import (
	"errors"
	//"fmt"
	"io"
	"math"
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	/**
	 * The default initial capacity for this table,
	 * used when not otherwise specified in a constructor.
	 */
	DEFAULT_INITIAL_CAPACITY int = 16

	/**
	 * The default load factor for this table, used when not
	 * otherwise specified in a constructor.
	 */
	DEFAULT_LOAD_FACTOR float32 = 0.75

	/**
	 * The default concurrency level for this table, used when not
	 * otherwise specified in a constructor.
	 */
	DEFAULT_CONCURRENCY_LEVEL int = 16

	/**
	 * The maximum capacity, used if a higher value is implicitly
	 * specified by either of the constructors with arguments.  MUST
	 * be a power of two <= 1<<30 to ensure that entries are indexable
	 * using ints.
	 */
	MAXIMUM_CAPACITY int = 1 << 30

	/**
	 * The maximum number of segments to allow; used to bound
	 * constructor arguments.
	 */
	MAX_SEGMENTS int = 1 << 16 // slightly conservative

	/**
	 * Number of unsynchronized retries in size and containsValue
	 * methods before resorting to locking. This is used to avoid
	 * unbounded retries if tables undergo continuous modification
	 * which would make it impossible to obtain an accurate result.
	 */
	RETRIES_BEFORE_LOCK int = 2
)

var (
	Debug           = false
	NilKeyError     = errors.New("Do not support nil as key")
	NilValueError   = errors.New("Do not support nil as value")
	NilActionError  = errors.New("Do not support nil as action")
	NonSupportKey   = errors.New("Non support for pointer, interface, channel, slice, map and function ")
	IllegalArgError = errors.New("IllegalArgumentException")
)

type Hashable interface {
	HashBytes() []byte
	Equals(v2 interface{}) bool
}

type hashEnginer struct {
	putFunc func(w io.Writer, v interface{})
}

//segments is read-only, don't need synchronized
type ConcurrentMap struct {
	engChecker *Once
	eng        unsafe.Pointer

	/**
	 * Kind of Reflect value for key
	 */
	kind unsafe.Pointer
	/**
	 * Mask value for indexing into segments. The upper bits of a
	 * key's hash code are used to choose the segment.
	 */
	segmentMask int

	/**
	 * Shift value for indexing within segments.
	 */
	segmentShift uint

	/**
	 * The segments, each of which is a specialized hash table
	 */
	segments []*Segment
}

/**
 * Returns the segment that should be used for key with given hash
 * @param hash the hash code for the key
 * @return the segment
 */
func (this *ConcurrentMap) segmentFor(hash uint32) *Segment {
	//默认segmentShift是28，segmentMask是（0xFFFFFFF）,hash>>this.segmentShift就是取前面4位
	//&segmentMask似乎没有必要
	//get first four bytes
	return this.segments[(hash>>this.segmentShift)&uint32(this.segmentMask)]
}

/**
 * Returns true if this map contains no key-value mappings.
 */
func (this *ConcurrentMap) IsEmpty() bool {
	segments := this.segments
	/*
	 * if any segment count isn't zero, Map will be no empty.
	 * 检查是否每个segment的count是否为0，并记录modCount和总和
	 */
	mc := make([]int32, len(segments))
	var mcsum int32 = 0
	for i := 0; i < len(segments); i++ {
		if atomic.LoadInt32(&segments[i].count) != 0 {
			return false
		} else {
			mc[i] = atomic.LoadInt32(&segments[i].modCount)
			mcsum += mc[i]
		}
	}

	/*
	 * if mcsum isn't zero, then modification is made,
	 * we will check per-segments count and if modCount be modified
	 * to avoid ABA problems in which an element in one segment was added and
	 * in another removed during traversal
	 */
	if mcsum != 0 {
		for i := 0; i < len(segments); i++ {
			if atomic.LoadInt32(&segments[i].count) != 0 || mc[i] != atomic.LoadInt32(&segments[i].modCount) {
				return false
			}
		}
	}
	return true
}

/**
 * Returns the number of key-value mappings in this map.
 */
func (this *ConcurrentMap) Size() int32 {
	segments := this.segments
	var sum int32 = 0
	var check int32 = 0
	mc := make([]int32, len(segments))

	// Try a few times to get accurate count. On failure due to
	// continuous async changes in table, resort to locking.
	for k := 0; k < RETRIES_BEFORE_LOCK; k++ {
		check = 0
		sum = 0
		var mcsum int32 = 0
		for i := 0; i < len(segments); i++ {
			sum += atomic.LoadInt32(&segments[i].count)
			mc[i] = atomic.LoadInt32(&segments[i].modCount)
			mcsum += mc[i]
		}
		if mcsum != 0 {
			for i := 0; i < len(segments); i++ {
				check += atomic.LoadInt32(&segments[i].count)
				if mc[i] != atomic.LoadInt32(&segments[i].modCount) {
					//async change happens, force retry
					check = -1 //
					break
				}
			}
		}

		//twoice counts ar same, it means no async change,
		//then will return sum
		if check == sum {
			break
		}
	}

	//async change happens in each loop
	//lock all segments to get accurate count
	if check != sum {
		sum = 0
		for i := 0; i < len(segments); i++ {
			segments[i].lock.Lock()
		}
		for i := 0; i < len(segments); i++ {
			sum += segments[i].count
		}
		for i := 0; i < len(segments); i++ {
			segments[i].lock.Unlock()
		}
	}
	return sum
}

/**
 * Returns the value to which the specified key is mapped,
 * or nil if this map contains no mapping for the key.
 */
func (this *ConcurrentMap) Get(key interface{}) (value interface{}, err error) {
	if isNil(key) {
		return nil, NilKeyError
	}
	//if atomic.LoadPointer(&this.kind) == nil {
	//	return nil, nil
	//}
	if hash, e := hashKey(key, this, false); e != nil {
		err = e
	} else {
		Printf("Get, %v, %v\n", key, hash)
		value = this.segmentFor(hash).get(key, hash)
	}
	return
}

/**
 * Tests if the specified object is a key in this table.
 *
 * @param  key   possible key
 * @return true if and only if the specified object is a key in this table,
 * as determined by the == method; false otherwise.
 */
func (this *ConcurrentMap) ContainsKey(key interface{}) (found bool, err error) {
	if isNil(key) {
		return false, NilKeyError
	}
	if atomic.LoadPointer(&this.kind) == nil {
		return false, nil
	}

	if hash, e := hashKey(key, this, false); e != nil {
		err = e
	} else {
		Printf("ContainsKey, %v, %v\n", key, hash)
		found = this.segmentFor(hash).containsKey(key, hash)
	}
	//hash := hash2(hashKey(key, this, false))
	//Printf("ContainsKey, %v, %v\n", key, hash)
	//found = this.segmentFor(hash).containsKey(key, hash)
	return
}

/**
 * Maps the specified key to the specified value in this table.
 * Neither the key nor the value can be nil.
 *
 * The value can be retrieved by calling the get method
 * with a key that is equal to the original key.
 *
 * @param key with which the specified value is to be associated
 * @param value to be associated with the specified key
 *
 * @return the previous value associated with key, or
 *         nil if there was no mapping for key
 */
func (this *ConcurrentMap) Put(key interface{}, value interface{}) (oldVal interface{}, err error) {
	if isNil(key) {
		return nil, NilKeyError
	}
	if isNil(value) {
		return nil, NilValueError
	}

	if hash, e := hashKey(key, this, false); e != nil {
		err = e
	} else {
		Printf("Put, %v, %v\n", key, hash)
		oldVal = this.segmentFor(hash).put(key, hash, value, false, nil)
	}
	//hash := hash2(hashKey(key, this, true))
	//Printf("Put, %v, %v\n", key, hash)
	//oldVal = this.segmentFor(hash).put(key, hash, value, false)
	return
}

/**
 * If mapping exists for the key, then maps the specified key to the specified value in this table.
 * else will ignore.
 * Neither the key nor the value can be nil.
 *
 * The value can be retrieved by calling the get method
 * with a key that is equal to the original key.
 *
 * @return the previous value associated with the specified key,
 *         or nil if there was no mapping for the key
 */
func (this *ConcurrentMap) PutIfAbsent(key interface{}, value interface{}) (oldVal interface{}, err error) {
	if isNil(key) {
		return nil, NilKeyError
	}
	if isNil(value) {
		return nil, NilValueError
	}

	if hash, e := hashKey(key, this, false); e != nil {
		err = e
	} else {
		Printf("PutIfAbsent, %v, %v\n", key, hash)
		oldVal = this.segmentFor(hash).put(key, hash, value, true, nil)
	}
	//hash := hash2(hashKey(key, this, true))
	//Printf("PutIfAbsent, %v, %v\n", key, hash)
	//oldVal = this.segmentFor(hash).put(key, hash, value, true)
	return
}

/**
 * Maps the specified key to the value that be returned by specified function in this table.
 * The key can not be nil.
 *
 * The value mapping specified key will be passed into action function as parameter.
 * If mapping does not exists for the key, nil will be passed into action function.
 * If return value by action function is nil, the specified key will be remove from map.
 *
 * @param key with which the specified value is to be associated
 * @param action that be called to generate new value mapping the specified key
 *
 * @return the previous value associated with key, or
 *         nil if there was no mapping for key
 */
func (this *ConcurrentMap) Update(key interface{}, action func(oldVal interface{}) (newVal interface{})) (oldVal interface{}, err error) {
	if isNil(key) {
		return nil, NilKeyError
	}
	if action == nil {
		return nil, NilActionError
	}

	if hash, e := hashKey(key, this, false); e != nil {
		err = e
	} else {
		Printf("Put, %v, %v\n", key, hash)
		oldVal = this.segmentFor(hash).put(key, hash, nil, false, action)
	}
	//hash := hash2(hashKey(key, this, true))
	//Printf("Put, %v, %v\n", key, hash)
	//oldVal = this.segmentFor(hash).put(key, hash, value, false)
	return
}

/**
 * Copies all of the mappings from the specified map to this one.
 * These mappings replace any mappings that this map had for any of the
 * keys currently in the specified map.
 *
 * @param m mappings to be stored in this map
 */
func (this *ConcurrentMap) PutAll(m map[interface{}]interface{}) (err error) {
	if isNil(m) {
		err = errors.New("Cannot copy nil map")
	}
	for k, v := range m {
		this.Put(k, v)
	}
	return
}

/**
 * Removes the key (and its corresponding value) from this map.
 * This method does nothing if the key is not in the map.
 *
 * @param  key the key that needs to be removed
 * @return the previous value associated with key, or nil if there was no mapping for key
 */
func (this *ConcurrentMap) Remove(key interface{}) (oldVal interface{}, err error) {
	if isNil(key) {
		return nil, NilKeyError
	}

	if hash, e := hashKey(key, this, false); e != nil {
		err = e
	} else {
		Printf("Remove, %v, %v\n", key, hash)
		oldVal = this.segmentFor(hash).remove(key, hash, nil)
	}
	//hash := hash2(hashKey(key, this, true))
	//Printf("Remove, %v, %v\n", key, hash)
	//oldVal = this.segmentFor(hash).remove(key, hash, nil)
	return
}

/**
 * Removes the mapping for the key and value from this map.
 * This method does nothing if no mapping for the key and value.
 *
 * @return true if mapping be removed, false otherwise
 */
func (this *ConcurrentMap) RemoveEntry(key interface{}, value interface{}) (ok bool, err error) {
	if isNil(key) {
		return false, NilKeyError
	}
	if isNil(value) {
		return false, NilValueError
	}

	if hash, e := hashKey(key, this, false); e != nil {
		err = e
	} else {
		Printf("RemoveEntry, %v, %v\n", key, hash)
		ok = this.segmentFor(hash).remove(key, hash, value) != nil
	}
	//hash := hash2(hashKey(key, this, true))
	//Printf("RemoveEntry, %v, %v\n", key, hash)
	//ok = this.segmentFor(hash).remove(key, hash, value) != nil
	return
}

/**
 * CompareAndReplace executes the compare-and-replace operation.
 * Replaces the value if the mapping exists for the previous and key from this map.
 * This method does nothing if no mapping for the key and value.
 *
 * @return true if value be replaced, false otherwise
 */
func (this *ConcurrentMap) CompareAndReplace(key interface{}, oldVal interface{}, newVal interface{}) (ok bool, err error) {
	if isNil(key) {
		return false, NilKeyError
	}
	if isNil(oldVal) || isNil(newVal) {
		return false, NilValueError
	}

	if hash, e := hashKey(key, this, false); e != nil {
		err = e
	} else {
		Printf("CompareAndReplace, %v, %v\n", key, hash)
		ok = this.segmentFor(hash).compareAndReplace(key, hash, oldVal, newVal)
	}
	//hash := hash2(hashKey(key, this, true))
	//Printf("CompareAndReplace, %v, %v\n", key, hash)
	//ok = this.segmentFor(hash).replaceWithOld(key, hash, oldVal, newVal)
	return
}

/**
 * Replaces the value if the key is in the map.
 * This method does nothing if no mapping for the key.
 *
 * @return the previous value associated with the specified key,
 *         or nil if there was no mapping for the key
 */
func (this *ConcurrentMap) Replace(key interface{}, value interface{}) (oldVal interface{}, err error) {
	if isNil(key) {
		return nil, NilKeyError
	}
	if isNil(value) {
		return nil, NilValueError
	}

	if hash, e := hashKey(key, this, false); e != nil {
		err = e
	} else {
		Printf("Replace, %v, %v\n", key, hash)
		oldVal = this.segmentFor(hash).replace(key, hash, value)
	}
	//hash := hash2(hashKey(key, this, true))
	//Printf("Replace, %v, %v\n", key, hash)
	//oldVal = this.segmentFor(hash).replace(key, hash, value)
	return
}

/**
 * Removes all of the mappings from this map.
 */
func (this *ConcurrentMap) Clear() {
	for i := 0; i < len(this.segments); i++ {
		this.segments[i].clear()
	}
}

//Iterator returns a iterator for ConcurrentMap
func (this *ConcurrentMap) Iterator() *MapIterator {
	return newMapIterator(this)
}

//ToSlice returns a slice that includes all key-value Entry in ConcurrentMap
func (this *ConcurrentMap) ToSlice() (kvs []*Entry) {
	kvs = make([]*Entry, 0, this.Size())
	itr := this.Iterator()
	for itr.HasNext() {
		kvs = append(kvs, itr.nextEntry())
	}
	return
}

func (this *ConcurrentMap) parseKey(key interface{}) (err error) {
	this.engChecker.Do(func() {
		var eng *hashEnginer

		val := key

		if _, ok := val.(Hashable); ok {
			eng = hasherEng
		} else {
			switch v := val.(type) {
			case bool:
				_ = v
				eng = boolEng
			case int:
				eng = intEng
			case int8:
				eng = int8Eng
			case int16:
				eng = int16Eng
			case int32:
				eng = int32Eng
			case int64:
				eng = int64Eng
			case uint:
				eng = uintEng
			case uint8:
				eng = uint8Eng
			case uint16:
				eng = uint16Eng
			case uint32:
				eng = uint32Eng
			case uint64:
				eng = uint64Eng
			case uintptr:
				eng = uintptrEng
			case float32:
				eng = float32Eng
			case float64:
				eng = float64Eng
			case complex64:
				eng = complex64Eng
			case complex128:
				eng = complex128Eng
			case string:
				eng = stringEng
			default:
				Printf("key = %v, other case\n", key)
				//some types can be used as key, we can use equals to test
				//_ = val == val

				rv := reflect.ValueOf(val)
				if ki, e := getKeyInfo(rv.Type()); e != nil {
					err = e
					return
				} else {
					putF := getPutFunc(ki)
					eng = &hashEnginer{}
					eng.putFunc = putF
				}
			}
		}

		this.eng = unsafe.Pointer(eng)

		Printf("key = %v, eng=%v, %v\n", key, this.eng, eng)
	})
	return
}

func (this *ConcurrentMap) newSegment(initialCapacity int, lf float32) (s *Segment) {
	s = new(Segment)
	s.loadFactor = lf
	table := make([]unsafe.Pointer, initialCapacity)
	s.setTable(table)
	s.lock = new(sync.Mutex)
	s.m = this
	return
}

func newConcurrentMap3(initialCapacity int,
	loadFactor float32, concurrencyLevel int) (m *ConcurrentMap) {
	m = &ConcurrentMap{}

	if !(loadFactor > 0) || initialCapacity < 0 || concurrencyLevel <= 0 {
		panic(IllegalArgError)
	}

	if concurrencyLevel > MAX_SEGMENTS {
		concurrencyLevel = MAX_SEGMENTS
	}

	// Find power-of-two sizes best matching arguments
	sshift := 0
	ssize := 1
	for ssize < concurrencyLevel {
		sshift++
		ssize = ssize << 1
	}

	m.segmentShift = uint(32) - uint(sshift)
	m.segmentMask = ssize - 1

	m.segments = make([]*Segment, ssize)

	if initialCapacity > MAXIMUM_CAPACITY {
		initialCapacity = MAXIMUM_CAPACITY
	}

	c := initialCapacity / ssize
	if c*ssize < initialCapacity {
		c++
	}
	cap := 1
	for cap < c {
		cap <<= 1
	}

	for i := 0; i < len(m.segments); i++ {
		m.segments[i] = m.newSegment(cap, loadFactor)
	}
	m.engChecker = new(Once)
	return
}

/**
 * Creates a new, empty map with the specified initial
 * capacity, load factor and concurrency level.
 *
 * @param initialCapacity the initial capacity. The implementation
 * performs internal sizing to accommodate this many elements.
 *
 * @param loadFactor  the load factor threshold, used to control resizing.
 * Resizing may be performed when the average number of elements per
 * bin exceeds this threshold.
 *
 * @param concurrencyLevel the estimated number of concurrently
 * updating threads. The implementation performs internal sizing
 * to try to accommodate this many threads.
 *
 * panic error "IllegalArgumentException" if the initial capacity is
 * negative or the load factor or concurrencyLevel are
 * nonpositive.
 *
 * Creates a new, empty map with a default initial capacity (16),
 * load factor (0.75) and concurrencyLevel (16).
 */
func NewConcurrentMap(paras ...interface{}) (m *ConcurrentMap) {
	ok := false
	cap := DEFAULT_INITIAL_CAPACITY
	factor := DEFAULT_LOAD_FACTOR
	concurrent_lvl := DEFAULT_CONCURRENCY_LEVEL

	if len(paras) >= 1 {
		if cap, ok = paras[0].(int); !ok {
			panic(IllegalArgError)
		}
	}

	if len(paras) >= 2 {
		if factor, ok = paras[1].(float32); !ok {
			panic(IllegalArgError)
		}
	}

	if len(paras) >= 3 {
		if concurrent_lvl, ok = paras[2].(int); !ok {
			panic(IllegalArgError)
		}
	}

	m = newConcurrentMap3(cap, factor, concurrent_lvl)
	return
}

/**
 * Creates a new map with the same mappings as the given map.
 * The map is created with a capacity of 1.5 times the number
 * of mappings in the given map or 16 (whichever is greater),
 * and a default load factor (0.75) and concurrencyLevel (16).
 *
 * @param m the map
 */
func NewConcurrentMapFromMap(m map[interface{}]interface{}) *ConcurrentMap {
	cm := newConcurrentMap3(int(math.Max(float64(float32(len(m))/DEFAULT_LOAD_FACTOR+1),
		float64(DEFAULT_INITIAL_CAPACITY))),
		DEFAULT_LOAD_FACTOR, DEFAULT_CONCURRENCY_LEVEL)
	cm.PutAll(m)
	return cm
}

/**
 * ConcurrentHashMap list entry.
 * Note only value field is variable and must use atomic to read/write it, other three fields are read-only after initializing.
 * so can use unsynchronized reader, the Segment.readValueUnderLock method is used as a
 * backup in case a nil (pre-initialized) value is ever seen in
 * an unsynchronized access method.
 */
type Entry struct {
	key   interface{}
	hash  uint32
	value unsafe.Pointer
	next  *Entry
}

func (this *Entry) Key() interface{} {
	return this.key
}

func (this *Entry) Value() interface{} {
	return *((*interface{})(atomic.LoadPointer(&this.value)))
}

func (this *Entry) fastValue() interface{} {
	return *((*interface{})(this.value))
}

func (this *Entry) storeValue(v *interface{}) {
	atomic.StorePointer(&this.value, unsafe.Pointer(v))
}

type Segment struct {
	m *ConcurrentMap //point to concurrentMap.eng, so it is **hashEnginer
	/**
	 * The number of elements in this segment's region.
	 * Must use atomic package's LoadInt32 and StoreInt32 functions to read/write this field
	 * otherwise read operation may cannot read latest value
	 */
	count int32

	/**
	 * Number of updates that alter the size of the table. This is
	 * used during bulk-read methods to make sure they see a
	 * consistent snapshot: If modCounts change during a traversal
	 * of segments computing size or checking containsValue, then
	 * we might have an inconsistent view of state so (usually)
	 * must retry.
	 */
	modCount int32

	/**
	 * The table is rehashed when its size exceeds this threshold.
	 * (The value of this field is always (int)(capacity *
	 * loadFactor).)
	 */
	threshold int32

	/**
	 * The per-segment table.
	 * Use unsafe.Pointer because must use atomic.LoadPointer function in read operations.
	 */
	pTable unsafe.Pointer //point to []unsafe.Pointer

	/**
	 * The load factor for the hash table. Even though this value
	 * is same for all segments, it is replicated to avoid needing
	 * links to outer object.
	 */
	loadFactor float32

	lock *sync.Mutex
}

func (this *Segment) enginer() *hashEnginer {
	return (*hashEnginer)(atomic.LoadPointer(&this.m.eng))
}

func (this *Segment) rehash() {
	oldTable := this.table() //*(*[]*Entry)(this.table)
	oldCapacity := len(oldTable)
	if oldCapacity >= MAXIMUM_CAPACITY {
		return
	}

	/*
	 * Reclassify nodes in each list to new Map.  Because we are
	 * using power-of-two expansion, the elements from each bin
	 * must either stay at same index, or move with a power of two
	 * offset. We eliminate unnecessary node creation by catching
	 * cases where old nodes can be reused because their next
	 * fields won't change. Statistically, at the default
	 * threshold, only about one-sixth of them need cloning when
	 * a table doubles. The nodes they replace will be garbage
	 * collectable as soon as they are no longer referenced by any
	 * reader thread that may be in the midst of traversing table
	 * right now.
	 */

	newTable := make([]unsafe.Pointer, oldCapacity<<1)
	atomic.StoreInt32(&this.threshold, int32(float32(len(newTable))*this.loadFactor))
	sizeMask := uint32(len(newTable) - 1)
	for i := 0; i < oldCapacity; i++ {
		// We need to guarantee that any existing reads of old Map can
		//  proceed. So we cannot yet nil out each bin.
		e := (*Entry)(oldTable[i])

		if e != nil {
			next := e.next
			//计算节点扩容后新的数组下标
			idx := e.hash & sizeMask

			//  Single node on list
			//如果没有后续的碰撞节点，直接复制到新数组即可
			if next == nil {
				newTable[idx] = unsafe.Pointer(e)
			} else {
				/* Reuse trailing consecutive sequence at same slot
				 * 数组扩容后原来数组下标相同（碰撞）的节点可能会计算出不同的新下标
				 * 如果把碰撞链表中所有节点的新下标列出，并将相邻的新下标相同的节点视为一段
				 * 那么下面的代码为了提高效率，会循环碰撞链表，找到链表中最后一段首节点（之后所有节点的新下标相同）
				 * 然后将这个首节点复制到新数组，后续节点因为计算出的新下标相同，所以在扩容后的数组中仍然在同一碰撞链表中
				 * 所以新的首节点的碰撞链表是正确的
				 * 新的首节点之外的其他现存碰撞链表上的节点，则重新复制到新节点（这个重要，可以保持旧节点的不变性）后放入新数组
				 * 这个过程的关键在于维持所有旧节点的next属性不会发生变化，这样才能让无锁的读操作保持线程安全
				 */
				lastRun := e
				lastIdx := idx
				for last := next; last != nil; last = last.next {
					k := last.hash & uint32(sizeMask)
					//发现新下标不同的节点就保存到lastIdx和lastRun中
					//所以lastIdx和lastRun总是对应现有碰撞链表中最后一段新下标相同节点的首节点和其对应的新下标
					if k != lastIdx {
						lastIdx = k
						lastRun = last
					}
				}
				newTable[lastIdx] = unsafe.Pointer(lastRun)

				// Clone all remaining nodes
				for p := e; p != lastRun; p = p.next {
					k := p.hash & sizeMask
					n := newTable[k]
					newTable[k] = unsafe.Pointer(&Entry{p.key, p.hash, p.value, (*Entry)(n)})
				}
			}
		}
	}
	atomic.StorePointer(&this.pTable, unsafe.Pointer(&newTable))
}

/**
 * Sets table to new pointer slice that all item points to HashEntry.
 * Call only while holding lock or in constructor.
 */
func (this *Segment) setTable(newTable []unsafe.Pointer) {
	this.threshold = (int32)(float32(len(newTable)) * this.loadFactor)
	this.pTable = unsafe.Pointer(&newTable)
}

/**
 * uses atomic to load table and returns.
 * Call while no lock.
 */
func (this *Segment) loadTable() (table []unsafe.Pointer) {
	return *(*[]unsafe.Pointer)(atomic.LoadPointer(&this.pTable))
}

/**
 * returns pointer slice that all item points to HashEntry.
 * Call only while holding lock or in constructor.
 */
func (this *Segment) table() []unsafe.Pointer {
	return *(*[]unsafe.Pointer)(this.pTable)
}

/**
 * Returns properly casted first entry of bin for given hash.
 */
func (this *Segment) getFirst(hash uint32) *Entry {
	tab := this.loadTable()
	return (*Entry)(atomic.LoadPointer(&tab[hash&uint32(len(tab)-1)]))
}

/**
 * Reads value field of an entry under lock. Called if value
 * field ever appears to be nil. see below code:
 * 		tab[index] = unsafe.Pointer(&Entry{key, hash, unsafe.Pointer(&value), first})
 * go memory model don't explain Entry initialization must be executed before
 * table assignment. So value is nil is possible only if a
 * compiler happens to reorder a HashEntry initialization with
 * its table assignment, which is legal under memory model
 * but is not known to ever occur.
 */
func (this *Segment) readValueUnderLock(e *Entry) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return e.fastValue()
}

/* Specialized implementations of map methods */

func (this *Segment) get(key interface{}, hash uint32) interface{} {
	if atomic.LoadInt32(&this.count) != 0 { // atomic-read
		e := this.getFirst(hash)
		for e != nil {
			if e.hash == hash && equals(e.key, key) {
				v := e.Value()
				if v != nil {
					//return
					return v
				}
				return this.readValueUnderLock(e) // recheck
			}
			e = e.next
		}
	}
	return nil
}

func (this *Segment) containsKey(key interface{}, hash uint32) bool {
	if atomic.LoadInt32(&this.count) != 0 { // read-volatile
		e := this.getFirst(hash)
		for e != nil {
			if e.hash == hash && equals(e.key, key) {
				return true
			}
			e = e.next
		}
	}
	return false
}

func (this *Segment) compareAndReplace(key interface{}, hash uint32, oldVal interface{}, newVal interface{}) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	e := this.getFirst(hash)
	for e != nil && (e.hash != hash || !equals(e.key, key)) {
		e = e.next
	}

	replaced := false
	if e != nil && oldVal == e.fastValue() {
		replaced = true
		e.storeValue(&newVal)
	}
	return replaced
}

func (this *Segment) replace(key interface{}, hash uint32, newVal interface{}) (oldVal interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	e := this.getFirst(hash)
	for e != nil && (e.hash != hash || !equals(e.key, key)) {
		e = e.next
	}

	if e != nil {
		oldVal = e.fastValue()
		e.storeValue(&newVal)
	}
	return
}

/**
 * put方法牵涉到count, modCount, pTable三个共享变量的修改
 * 在Java中count和pTable是volatile字段，而modCount不是
 * 由于IsEmpty和Size等操作会读取count, modCount和pTable并且是无锁的，这里有必要对进行并发安全性的分析
 * 在Java中，volatile的读具有Acquire语义，volatile的写具有release语义，而put的最后会写入count，
 * 其他读操作总是会先读取count，由此保证了put中其他的写入操作不会被reorder到写入count之后，而读操作中其他的读取不会被reorder到读count之前
 * 由此保证了多线程情况下读和写线程中看到的操作次序不会发送混乱，
 * 在Golang中，StorePointer内部使用了xchgl指令，具有内存屏障，但是Load操作似乎并未具有明确的acquire语义
 */
func (this *Segment) put(key interface{}, hash uint32, value interface{}, onlyIfAbsent bool, action func(oldValue interface{}) (newVal interface{})) (oldValue interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()

	c := this.count
	if c > this.threshold { // ensure capacity
		this.rehash()
	}

	tab := this.table()
	index := hash & uint32(len(tab)-1)
	first := (*Entry)(tab[index])
	e := first

	for e != nil && (e.hash != hash || !equals(e.key, key)) {
		e = e.next
	}

	if action == nil {
		if e != nil {
			oldValue = e.fastValue()
			if !onlyIfAbsent {
				e.storeValue(&value)
			}
		} else {
			c++
			oldValue = nil
			this.modCount++
			tab[index] = unsafe.Pointer(&Entry{key, hash, unsafe.Pointer(&value), first})
			atomic.StoreInt32(&this.count, c) // atomic write 这里可以保证对modCount和tab的修改不会被reorder到this.count之后
		}
	} else {
		if e != nil {
			oldValue = e.fastValue()
		} else {
			c++
			oldValue = nil
		}

		newVal := action(oldValue)
		if newVal != nil {
			if oldValue == nil {
				e = &Entry{key, hash, unsafe.Pointer(&value), first}
				tab[index] = unsafe.Pointer(e)
				this.modCount++
				atomic.StoreInt32(&this.count, c) // atomic write 这里可以保证对modCount和tab的修改不会被reorder到this.count之后
			}
			e.storeValue(&newVal)
		} else if e != nil {
			//remove key if action returns nil
			c--
			this.modCount++
			newFirst := e.next
			for p := first; p != e; p = p.next {
				newFirst = &Entry{p.key, p.hash, p.value, newFirst}
			}
			tab[index] = unsafe.Pointer(newFirst)
			atomic.StoreInt32(&this.count, c) //this.count = c
		}
	}
	return
}

/**
 * Remove; match on key only if value nil, else match both.
 */
func (this *Segment) remove(key interface{}, hash uint32, value interface{}) (oldValue interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()

	c := this.count - 1
	tab := this.table()
	index := hash & uint32(len(tab)-1)
	first := (*Entry)(tab[index])
	e := first

	for e != nil && (e.hash != hash || !equals(e.key, key)) {
		e = e.next
	}

	if e != nil {
		v := e.fastValue()
		if value == nil || value == v {
			oldValue = v
			// All entries following removed node can stay
			// in list, but all preceding ones need to be
			// cloned.
			this.modCount++
			newFirst := e.next
			for p := first; p != e; p = p.next {
				newFirst = &Entry{p.key, p.hash, p.value, newFirst}
			}
			tab[index] = unsafe.Pointer(newFirst)
			atomic.StoreInt32(&this.count, c) //this.count = c
		}
	}
	return
}

func (this *Segment) clear() {
	if atomic.LoadInt32(&this.count) != 0 {
		this.lock.Lock()
		defer this.lock.Unlock()

		tab := this.table()
		for i := 0; i < len(tab); i++ {
			tab[i] = nil
		}
		this.modCount++
		atomic.StoreInt32(&this.count, 0) //this.count = 0 // write-volatile
	}
}

/**
 * Applies a supplemental hash function to a given hashCode, which
 * defends against poor quality hash functions.  This is critical
 * because ConcurrentHashMap uses power-of-two length hash tables,
 * that otherwise encounter collisions for hashCodes that do not
 * differ in lower or upper bits.
 */
func hash2(h uint32) uint32 {
	//// Spread bits to regularize both segment and index locations,
	//// using variant of single-word Wang/Jenkins hash.
	//h += (h << 15) ^ 0xffffcd7d
	//h ^= (h >> 10)
	//h += (h << 3)
	//h ^= (h >> 6)
	//h += (h << 2) + (h << 14)
	//return uint32(h ^ (h >> 16))

	//Now all hashcode is created by FNVa, it isn't a poor quality hash function
	//so I removes the hash operation for second time
	return h
}

/* ---------------- Iterator Support -------------- */

type MapIterator struct {
	nextSegmentIndex int
	nextTableIndex   int
	currentTable     []unsafe.Pointer
	nextE            *Entry
	lastReturned     *Entry
	cm               *ConcurrentMap
}

func (this *MapIterator) advance() {
	if this.nextE != nil {
		this.nextE = this.nextE.next
		if this.nextE != nil {
			return
		}
	}

	for this.nextTableIndex >= 0 {
		this.nextE = (*Entry)(atomic.LoadPointer(&this.currentTable[this.nextTableIndex]))
		this.nextTableIndex--
		if this.nextE != nil {
			return
		}
	}

	for this.nextSegmentIndex >= 0 {
		seg := this.cm.segments[this.nextSegmentIndex]
		this.nextSegmentIndex--
		if atomic.LoadInt32(&seg.count) != 0 {
			this.currentTable = seg.loadTable()
			for j := len(this.currentTable) - 1; j >= 0; j-- {
				this.nextE = (*Entry)(atomic.LoadPointer(&this.currentTable[j]))
				if this.nextE != nil {
					this.nextTableIndex = j - 1
					return
				}
			}
		}
	}
}

func (this *MapIterator) HasNext() bool {
	return this.nextE != nil
}

func (this *MapIterator) Next() (key interface{}, value interface{}, ok bool) {
	if this.nextE == nil {
		return nil, nil, false
	}
	this.lastReturned = this.nextE
	this.advance()
	key, value, ok = this.lastReturned.Key(), this.lastReturned.Value(), true
	return
}

func (this *MapIterator) Remove() (ok bool) {
	if this.lastReturned == nil {
		return false
	}
	this.cm.Remove(this.lastReturned.key)
	this.lastReturned = nil
	return true
}

func (this *MapIterator) nextEntry() *Entry {
	if this.nextE == nil {
		panic("IllegalStateException")
	}
	this.lastReturned = this.nextE
	this.advance()
	return this.lastReturned
}

func newMapIterator(cm *ConcurrentMap) *MapIterator {
	hi := MapIterator{}
	hi.nextSegmentIndex = len(cm.segments) - 1
	hi.nextTableIndex = -1
	hi.cm = cm
	hi.advance()
	return &hi
}
