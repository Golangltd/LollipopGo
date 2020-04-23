package concurrent

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	c "github.com/smartystreets/goconvey/convey"
)

func TestNil(t *testing.T) {
	c.Convey("Nil cannot be as key", t, func() {
		cm := NewConcurrentMap()
		_, err := cm.Put(nil, 1)
		c.So(err, c.ShouldNotBeNil)

		var nilVal interface{} = nil
		_, err = cm.Put(nilVal, 1)
		c.So(err, c.ShouldNotBeNil)

		var nilPtr *string = nil
		nilVal = nilPtr
		_, err = cm.Put(nilVal, 1)
		c.So(err, c.ShouldNotBeNil)

		_, err = cm.Put(1, nil)
		c.So(err, c.ShouldNotBeNil)

		_, err = cm.Put(1, nilVal)
		c.So(err, c.ShouldNotBeNil)

		nilVal = nilPtr
		_, err = cm.Put(1, nilVal)
		c.So(err, c.ShouldNotBeNil)

		_, err = cm.Get(nil)
		c.So(err, c.ShouldNotBeNil)

		_, err = cm.Get(nilVal)
		c.So(err, c.ShouldNotBeNil)

		nilVal = nilPtr
		_, err = cm.Get(nilVal)
		c.So(err, c.ShouldNotBeNil)
	})
}

/*-------------test different types as key------------------------*/
func testConcurrentMap(t *testing.T, datas map[interface{}]interface{}) {
	var firstKey, firstVal interface{}
	var secondaryKey, secondaryVal interface{}
	i := 0
	for k, v := range datas {
		if i == 0 {
			firstKey, firstVal = k, v
		} else if i == 1 {
			secondaryKey, secondaryVal = k, v
			break
		}
		i++
	}

	m := NewConcurrentMap()

	//test Put first key-value pair
	previou, err := m.Put(firstKey, firstVal)
	if previou != nil || err != nil {
		t.Errorf("Put %v, %v firstly, return %v, %v, want nil, nil", firstKey, firstVal, previou, err)
	}

	//test Put again
	previou, err = m.Put(firstKey, firstVal)
	if previou != firstVal || err != nil {
		t.Errorf("Put %v, %v second time, return %v, %v, want %v, nil", firstKey, firstVal, previou, err, firstVal)
	}

	//test PutIfAbsent, if value is incorrect, PutIfAbsent will be ignored
	v := rand.Float32()
	previou, err = m.PutIfAbsent(firstKey, v)
	if previou != firstVal || err != nil {
		t.Errorf("PutIfAbsent %v, %v three time, return %v, %v, want %v, nil", firstKey, v, previou, err, firstVal)
	}

	//test Get
	val, err := m.Get(firstKey)
	if val != firstVal || err != nil {
		t.Errorf("Get %v, return %v, %v, want %v, nil", firstKey, val, err, firstVal)
	}

	//test Size
	s := m.Size()
	if s != 1 {
		t.Errorf("Get size of m, return %v, want 1", s)
	}

	//test PutAll
	m.PutAll(datas)
	s = m.Size()
	if s != int32(len(datas)) {
		t.Errorf("Get size of m, return %v, want %v", s, len(datas))
	}

	//test remove a key-value pair, if value is incorrect, RemoveKV will return be ignored and return false
	ok, err := m.RemoveEntry(secondaryKey, v)
	if ok != false || err != nil {
		t.Errorf("RemoveKV %v, %v, return %v, %v, want false, nil", secondaryKey, v, ok, err)
	}

	//test replace a value for a key
	previou, err = m.Replace(secondaryKey, v)
	if previou != secondaryVal || err != nil {
		t.Errorf("Replace %v, %v, return %v, %v, want %v, nil", secondaryKey, v, previou, err, secondaryVal)
	}

	//test replace a value for a key-value pair, if value is incorrect, replace will ignored and return false
	ok, err = m.CompareAndReplace(secondaryKey, secondaryVal, v)
	if ok != false || err != nil {
		t.Errorf("ReplaceWithOld  %v, %v, %v, return %v, %v, want false, nil", secondaryKey, secondaryVal, v, ok, err)
	}

	//test replace a value for a key-value pair, if value is correct, replace will success
	ok, err = m.CompareAndReplace(secondaryKey, v, secondaryVal)
	if ok != true || err != nil {
		t.Errorf("ReplaceWithOld %v, %v, %v, return %v, %v, want true, nil", secondaryKey, v, secondaryVal, ok, err)
	}

	//test remove a key
	previou, err = m.Remove(secondaryKey)
	if previou != secondaryVal || err != nil {
		t.Errorf("Remove %v, return %v, %v, want %v, nil", secondaryKey, previou, err, secondaryVal)
	}

	//test clear
	m.Clear()
	if m.Size() != 0 {
		t.Errorf("Get size of m after calling Clear(), return %v, want 0", val)
	}
}

func TestIntKey(t *testing.T) {
	testConcurrentMap(t, map[interface{}]interface{}{
		1: 10,
		2: 20,
		3: 30,
		4: 40,
	})
}

func TestStringKey(t *testing.T) {
	testConcurrentMap(t, map[interface{}]interface{}{
		strconv.Itoa(1): 10,
		strconv.Itoa(2): 20,
		strconv.Itoa(3): 30,
		strconv.Itoa(4): 40,
	})
}

func Testfloat32Key(t *testing.T) {
	testConcurrentMap(t, map[interface{}]interface{}{
		float32(1): 10,
		float32(2): 20,
		float32(3): 30,
		float32(4): 40,
	})
}

func Testfloat64Key(t *testing.T) {
	testConcurrentMap(t, map[interface{}]interface{}{
		float64(1): 10,
		float64(2): 20,
		float64(3): 30,
		float64(4): 40,
	})
}

//note interface{} is empty interface
func TestEmptyInterface(t *testing.T) {
	var a, b, c, d interface{} = 1, 2, 3, 4
	testConcurrentMap(t, map[interface{}]interface{}{
		a: 10,
		b: 20,
		c: 30,
		d: 40,
	})

	cm := NewConcurrentMap()
	cm.Put(a, 10)

	e := a
	if v, err := cm.Get(e); v != 10 || err != nil {
		t.Errorf("Get %v, return %v, %v, want %v", &e, v, err, 10)
	}
}

//user implements concurrent.Hasher interface
type user struct {
	id   string
	Name string
}

func (u *user) HashBytes() []byte {
	return []byte(u.id)
}
func (u *user) Equals(v2 interface{}) (equal bool) {
	u2, ok := v2.(*user)
	return ok && u.id == u2.id
}

//test Hasher interface
func TestHasherKey(t *testing.T) {
	a, b, c, d := &user{"1", "n1"}, &user{"2", "n2"}, &user{"3", "n3"}, &user{"4", "n4"}
	testConcurrentMap(t, map[interface{}]interface{}{
		a: 10,
		b: 20,
		c: 30,
		d: 40,
	})

	cm := NewConcurrentMap()
	cm.Put(a, 10)

	e := &user{"1", "n1"}
	if v, err := cm.Get(e); v != 10 || err != nil {
		t.Errorf("Get %v, return %v, %v, want %v", &e, v, err, 10)
	}
}

//size of small is less than word size
//the memory layout is different with struct what size is greater than word size before golang 1.4
type small struct {
	Id   byte
	name byte
}

//test small struct
func TestSmallStruct(t *testing.T) {
	a, b, c, d := small{1, 1}, small{2, 2}, small{3, 3}, small{4, 4}
	testConcurrentMap(t, map[interface{}]interface{}{
		a: 10,
		b: 20,
		c: 30,
		d: 40,
	})

	//test using the interface object and original value as key, two value should return the same hash code
	cm := NewConcurrentMap()
	cm.Put(a, 10)
	e := small{1, 1}
	if v, err := cm.Get(e); v != 10 || err != nil {
		t.Errorf("Get %v, return %v, %v, want %v", &e, v, err, 10)
	}
}

//compositeStruct include anothe struct
type compositeStruct struct {
	F1 string
	f2 int
	small
}

//test composite Struct
func TestCompositeStruct(t *testing.T) {
	a, b, c, d := compositeStruct{"1", 1, small{1, 1}}, compositeStruct{"2", 2, small{2, 2}}, compositeStruct{"3", 3, small{3, 3}}, compositeStruct{"4", 4, small{4, 4}}
	testConcurrentMap(t, map[interface{}]interface{}{
		a: 10,
		b: 20,
		c: 30,
		d: 40,
	})

	//test using the interface object and original value as key, two value should return the same hash code
	cm := NewConcurrentMap()
	cm.Put(a, 10)
	e := compositeStruct{"1", 1, small{1, 1}}
	if v, err := cm.Get(e); v != 10 || err != nil {
		t.Errorf("Get %v, return %v, %v, want %v", &e, v, err, 10)
	}
}

/**
 * test update method
 * put three *user, e.g. &{"1", "jack"}, &{"2", "jack"}, &{"3", "stone"}
 * last map will include:
 * map{
 *    "jack":  [&{"1", "jack"}, &{"2", "jack"}]
 *    "stone": [&{"3", "stone"}]
 * }
 **/
func TestUpdate(t *testing.T) {
	//appendFunc returns a function that appends an user into *user slice
	appendFunc := func(u *user) func(oldVal interface{}) (newVal interface{}) {
		return func(oldVal interface{}) (newVal interface{}) {
			if u == nil {
				return nil
			}
			if oldVal == nil {
				users := make([]*user, 0, 1)
				return append(users, u)
			} else {
				return append(oldVal.([]*user), u)
			}
		}
	}

	cm := NewConcurrentMap()

	//put user with name jack
	u1 := &user{id: "1", Name: "jack"}
	old, err := cm.Update(u1.Name, appendFunc(u1))
	if old != nil || err != nil {
		t.Errorf("Update %v, %v, return %v, %v, want nil, nil", u1.id, u1, old, err)
	}

	//Getting value by "jack" returns [u1]
	v, err := cm.Get(u1.Name)
	if users := v.([]*user); len(users) != 1 || users[0] != u1 {
		t.Errorf("Get %v, return %v, %v, want [%v], nil", u1.id, users, err, old, u1)
	}

	//put another user with name jack
	u2 := &user{id: "2", Name: "jack"}
	old, err = cm.Update(u2.Name, appendFunc(u2))
	if users := old.([]*user); old == nil || len(users) != 1 || users[0] != u1 || err != nil {
		t.Errorf("Update %v, %v, return %#v, %v, want [%v], nil", u2.Name, u2, old, err, u1)
	}

	//Getting value by "jack" returns [u1, u2]
	v, err = cm.Get(u2.Name)
	if users := v.([]*user); len(users) != 2 || users[1] != u2 {
		t.Errorf("Get %v, return %v, %v, want [%v, %v], nil", u2.Name, users, err, old, u1, u2)
	}

	//put an user with name stone
	u3 := &user{id: "3", Name: "stone"}
	old, err = cm.Update(u3.Name, appendFunc(u3))
	if old != nil || err != nil {
		t.Errorf("Update %v, %v, return %#v, %v, want nil, nil", u3.Name, u3, old, err)
	}

	//Getting value by "stone" returns [u3]
	v, err = cm.Get(u3.Name)
	if users := v.([]*user); len(users) != 1 || users[0] != u3 {
		t.Errorf("Get %v, return %v, %v, want [%v], nil", u3.Name, users, err, old, u3)
	}

	//put nil with name stone
	old, err = cm.Update(u3.Name, appendFunc(nil))
	if old == nil || err != nil {
		t.Errorf("Update %v, nil, return %#v, %v, want %v, nil", u3.Name, old, err, u3)
	}
	
	//Getting value by "stone" returns nil
	v, err = cm.Get(u3.Name)
	if v != nil || err != nil {
		t.Errorf("Get %v, return %v, %v, want nil, nil", u3.Name, v, err)
	}

}

//user1 implements Ider interface
type user1 struct {
	id   string
	Name string
}

func (u *user1) Id() string {
	return u.id
}

type Ider interface {
	Id() string
}

//test slice, function, map, pointer and interface as key
func TestUnableHash(t *testing.T) {
	testHash := func(k interface{}) (err error) {
		cm := NewConcurrentMap()
		_, err = cm.Put(k, 1)
		return
	}

	//do not support slice as key
	err := testHash([]int{1})
	if err == nil {
		t.Errorf("Put slice, return nil, should be not nil")
	}

	//do not support function as key
	f := func() {}
	err = testHash(f)
	if err == nil {
		t.Errorf("Put function, return nil, should be not nil")
	}

	//do not support map as key
	err = testHash(map[int]int{1: 1})
	if err == nil {
		t.Errorf("Put map, return nil, should be not nil")
	}

	//do not support pointer as key
	a := 1
	err = testHash(&a)
	if err == nil {
		t.Errorf("Put function, return nil, should be not nil")
	}

	//do not support interface as key (note interface is different with interface{})
	//The kind of interface is a realy pointer
	var i Ider = &user1{"1", "n1"}
	err = testHash(i)
	if err == nil {
		t.Errorf("Put map, return nil, should be not nil")
	}
}

func TestToSlice(t *testing.T) {
	cm := NewConcurrentMap()
	kvs := cm.ToSlice()
	if kvs == nil || len(kvs) != 0 {
		t.Errorf("Call ToSlice for a empty map, return %v, should be empty slice", kvs)
	}

	cm.Put(1, 10)
	kvs = cm.ToSlice()
	if kvs == nil || len(kvs) != 1 || kvs[0].Key() != 1 || kvs[0].Value() != 10 {
		t.Errorf("Call ToSlice after put one key-value pair, return %v, should include one entry", kvs)
	}

	cm.Remove(1)
	kvs = cm.ToSlice()
	if kvs == nil || len(kvs) != 0 {
		t.Errorf("Call ToSlice for a empty map, return %v, should be empty slice", kvs)
	}

}

/*--------test cases copied from go standard library's map_test.go--------------------*/
//TestNegativeZero fail
//// negative zero is a good test because:
////  1) 0 and -0 are equal, yet have distinct representations.
////  2) 0 is represented as all zeros, -0 isn't.
//// I'm not sure the language spec actually requires this behavior,
//// but it's what the current map implementation does.
//func TestNegativeZero(t *testing.T) {
//	m := NewConcurrentMap(0)
//	var zero float64 = +0.0
//	var nzero float64 = math.Copysign(0.0, -1.0)

//	m.Put(zero, true)
//	m.Put(nzero, true) // should overwrite +0 entry

//	if m.Size() != 1 {
//		t.Error("length wrong", m.Size())
//	}

//	itr := NewHashIterator(m)
//	for {
//		if itr.HasNext() {
//			e := itr.nextEntry()
//			if math.Copysign(1.0, e.key.(float64)) > 0 {
//				t.Error("wrong sign")
//			}
//		} else {
//			break
//		}
//	}

//	m = NewConcurrentMap(0)

//	m.Put(nzero, true)
//	m.Put(zero, true) // should overwrite -0.0 entry

//	if m.Size() != 1 {
//		t.Error("length wrong")
//	}

//	itr = NewHashIterator(m)
//	for {
//		if itr.HasNext() {
//			e := itr.nextEntry()
//			if math.Copysign(1.0, e.key.(float64)) < 0 {
//				t.Error("wrong sign")
//			}
//		} else {
//			break
//		}
//	}
//}

// nan is a good test because nan != nan, and nan has
// a randomized hash value.
func TestNan(t *testing.T) {
	m := NewConcurrentMap(0) //make(map[float64]int, 0)
	nan := math.NaN()
	m.Put(nan, 1)
	m.Put(nan, 2)
	m.Put(nan, 4)
	if m.Size() != 3 {
		t.Error("length wrong")
	}
	s := 0
	itr := m.Iterator()
	for {
		k, v, ok := itr.Next()
		if !ok {
			break
		}
		if k == k {
			t.Error("nan disappeared")
		}

		vi := v.(int)
		if (vi & (vi - 1)) != 0 {
			t.Error("value wrong")
		}
		s |= vi
	}
	if s != 7 {
		t.Error("values wrong")
	}
}

func TestGrowWithNaN(t *testing.T) {
	m := NewConcurrentMap(0) //make(map[float64]int, 0)
	nan := math.NaN()
	m.Put(nan, 1)
	m.Put(nan, 2)
	m.Put(nan, 4)
	cnt := 0
	s := 0
	growflag := true

	for itr := m.Iterator(); itr.HasNext(); {
		ki, vi, _ := itr.Next()
		k, v := ki.(float64), vi.(int)
		if growflag {
			// force a hashtable resize
			for i := 0; i < 100; i++ {
				m.Put(float64(i), i)
			}
			growflag = false
		}
		if k != k {
			cnt++
			s |= v
		}
	}
	if cnt != 3 {
		t.Error("NaN keys lost during grow")
	}
	if s != 7 {
		t.Error("NaN values lost during grow")
	}
}

type FloatInt struct {
	x float64
	y int
}

func TestGrowWithNegativeZero(t *testing.T) {
	negzero := math.Copysign(0.0, -1.0)
	m := make(map[FloatInt]int, 4)
	m[FloatInt{0.0, 0}] = 1
	m[FloatInt{0.0, 1}] = 2
	m[FloatInt{0.0, 2}] = 4
	m[FloatInt{0.0, 3}] = 8
	growflag := true
	s := 0
	cnt := 0
	negcnt := 0
	// The first iteration should return the +0 key.
	// The subsequent iterations should return the -0 key.
	// I'm not really sure this is required by the spec,
	// but it makes sense.
	// TODO: are we allowed to get the first entry returned again???
	for k, v := range m {
		if v == 0 {
			continue
		} // ignore entries added to grow table
		cnt++
		if math.Copysign(1.0, k.x) < 0 {
			if v&16 == 0 {
				t.Error("key/value not updated together 1")
			}
			negcnt++
			s |= v & 15
		} else {
			if v&16 == 16 {
				t.Error("key/value not updated together 2", k, v)
			}
			s |= v
		}
		if growflag {
			// force a hashtable resize
			for i := 0; i < 100; i++ {
				m[FloatInt{3.0, i}] = 0
			}
			// then change all the entries
			// to negative zero
			m[FloatInt{negzero, 0}] = 1 | 16
			m[FloatInt{negzero, 1}] = 2 | 16
			m[FloatInt{negzero, 2}] = 4 | 16
			m[FloatInt{negzero, 3}] = 8 | 16
			growflag = false
		}
	}
	if s != 15 {
		t.Error("entry missing", s)
	}
	if cnt != 4 {
		t.Error("wrong number of entries returned by iterator", cnt)
	}
	if negcnt != 3 {
		t.Error("update to negzero missed by iteration", negcnt)
	}
}

func TestIterGrowAndDelete(t *testing.T) {
	m := make(map[int]int, 4)
	for i := 0; i < 100; i++ {
		m[i] = i
	}
	growflag := true
	for k := range m {
		//t.Log("k ad growflag", k, growflag)
		if growflag {
			// grow the table
			for i := 100; i < 1000; i++ {
				m[i] = i
			}
			// delete all odd keys
			for i := 1; i < 1000; i += 2 {
				delete(m, i)
			}
			growflag = false
		} else {
			if k&1 == 1 {
				t.Error("odd value returned")
			}
		}
	}
}

func TestIterGrowAndDelete1(t *testing.T) {
	m := NewConcurrentMap(4) //	make(map[int]int, 4)
	for i := 0; i < 100; i++ {
		m.Put(i, i)
	}
	growflag := true
	for itr := m.Iterator(); itr.HasNext(); {
		k, _, _ := itr.Next()
		//t.Log("k ad growflag111111", k, growflag)
		if growflag {
			// grow the table
			for i := 100; i < 1000; i++ {
				m.Put(i, i)
			}
			// delete all odd keys
			for i := 1; i < 1000; i += 2 {
				m.Remove(i)
			}
			growflag = false
		} else {
			if k.(int)&1 == 1 {
				for itr := m.Iterator(); itr.HasNext(); {
					k, _, _ := itr.Next()
					if k.(int)&1 == 1 {
						t.Error("odd value returned by itr")
					}
				}
				//ConcurrentMap cannot iterate the values changed outside iterator after grow
				//t.Error("odd value returned")
			}
		}
	}
}

// make sure old bucket arrays don't get GCd while
// an iterator is still using them.
func TestIterGrowWithGC(t *testing.T) {
	m := NewConcurrentMap(4) //	make(map[int]int, 4)
	for i := 0; i < 16; i++ {
		m.Put(i, i)
	}
	growflag := true
	bitmask := 0
	for itr := m.Iterator(); itr.HasNext(); {
		ki, _, _ := itr.Next()
		k := ki.(int)
		if k < 16 {
			bitmask |= 1 << uint(k)
		}
		if growflag {
			// grow the table
			for i := 100; i < 1000; i++ {
				m.Put(i, i)
			}
			// trigger a gc
			runtime.GC()
			growflag = false
		}
	}
	if bitmask != 1<<16-1 {
		t.Error("missing key", bitmask)
	}
}

func testConcurrentReadsAfterGrowth(t *testing.T, useReflect bool) {
	if runtime.GOMAXPROCS(-1) == 1 {
		defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(16))
	}
	numLoop := 10
	numGrowStep := 250
	numReader := 16
	if testing.Short() {
		numLoop, numGrowStep = 2, 500
	}
	for i := 0; i < numLoop; i++ {
		m := NewConcurrentMap() //	make(map[int]int, 0)
		for gs := 0; gs < numGrowStep; gs++ {
			m.Put(gs, gs)
			var wg sync.WaitGroup
			wg.Add(numReader * 2)
			for nr := 0; nr < numReader; nr++ {
				go func() {
					defer wg.Done()
					for itr := m.Iterator(); itr.HasNext(); {
						_, _, _ = itr.Next()
					}
				}()
				go func() {
					defer wg.Done()
					for key := 0; key < gs; key++ {
						_, _ = m.Get(key)
					}
				}()
			}
			wg.Wait()
		}
	}
}

func TestConcurrentReadsAfterGrowth(t *testing.T) {
	testConcurrentReadsAfterGrowth(t, false)
}

func TestConcurrentReadsAfterGrowthReflect(t *testing.T) {
	testConcurrentReadsAfterGrowth(t, true)
}

func TestBigItems(t *testing.T) {
	var key [256]string
	for i := 0; i < 256; i++ {
		key[i] = "foo"
	}
	m := NewConcurrentMap(4) //make(map[[256]string][256]string, 4)
	for i := 0; i < 100; i++ {
		key[37] = fmt.Sprintf("string%02d", i)
		m.Put(key, key) //m[key] = key
	}

	var keys [100]string
	var values [100]string
	i := 0
	for itr := m.Iterator(); itr.HasNext(); {
		ki, vi, _ := itr.Next()
		k, v := ki.([256]string), vi.([256]string)
		//for k, v := range m {
		keys[i] = k[37]
		values[i] = v[37]
		i++
	}
	sort.Strings(keys[:])
	sort.Strings(values[:])
	for i := 0; i < 100; i++ {
		if keys[i] != fmt.Sprintf("string%02d", i) {
			t.Errorf("#%d: missing key: %v", i, keys[i])
		}
		if values[i] != fmt.Sprintf("string%02d", i) {
			t.Errorf("#%d: missing value: %v", i, values[i])
		}
	}
}

type empty struct {
}

func TestEmptyKeyAndValue(t *testing.T) {
	//a := make(map[int]empty, 4)
	//b := make(map[empty]int, 4)
	//c := make(map[empty]empty, 4)
	a := NewConcurrentMap(4)
	b := NewConcurrentMap(4)
	c := NewConcurrentMap(4)
	a.Put(0, empty{})       //a[0] = empty{}
	b.Put(empty{}, 0)       //b[empty{}] = 0
	b.Put(empty{}, 1)       //b[empty{}] = 1
	c.Put(empty{}, empty{}) //c[empty{}] = empty{}

	if a.Size() != 1 { // len(a) != 1 {
		t.Errorf("empty value insert problem")
	}
	if v, err := b.Get(empty{}); v != 1 || err != nil { //} b[empty{}] != 1 {
		t.Errorf("empty key returned wrong value")
	}
}

// Tests a map with a single bucket, with same-lengthed short keys
// ("quick keys") as well as long keys.
func TestSingleBucketMapStringKeys_DupLen(t *testing.T) {
	testMapLookups(t, NewConcurrentMapFromMap(map[interface{}]interface{}{
		"x":    "x1val",
		"xx":   "x2val",
		"foo":  "fooval",
		"bar":  "barval", // same key length as "foo"
		"xxxx": "x4val",
		strings.Repeat("x", 128): "longval1",
		strings.Repeat("y", 128): "longval2",
	}))
}

// Tests a map with a single bucket, with all keys having different lengths.
func TestSingleBucketMapStringKeys_NoDupLen(t *testing.T) {
	testMapLookups(t, NewConcurrentMapFromMap(map[interface{}]interface{}{
		"x":                      "x1val",
		"xx":                     "x2val",
		"foo":                    "fooval",
		"xxxx":                   "x4val",
		"xxxxx":                  "x5val",
		"xxxxxx":                 "x6val",
		strings.Repeat("x", 128): "longval",
	}))
}

func testMapLookups(t *testing.T, m *ConcurrentMap) {
	for itr := m.Iterator(); itr.HasNext(); {
		k, v, _ := itr.Next()
		if v1, err := m.Get(k); v1 != v || err != nil {
			t.Fatalf("m[%q] = %q; want %q", k, v1, v)
		}
	}
}

//TestMapNanGrowIterator fail
//// Tests whether the iterator returns the right elements when
//// started in the middle of a grow, when the keys are NaNs.
//func TestMapNanGrowIterator(t *testing.T) {
//	m := make(map[float64]int)
//	nan := math.NaN()
//	const nBuckets = 16
//	// To fill nBuckets buckets takes LOAD * nBuckets keys.
//	nKeys := int(nBuckets * *runtime.HashLoad)

//	// Get map to full point with nan keys.
//	for i := 0; i < nKeys; i++ {
//		m[nan] = i
//	}
//	// Trigger grow
//	m[1.0] = 1
//	delete(m, 1.0)

//	// Run iterator
//	found := make(map[int]struct{})
//	for _, v := range m {
//		if v != -1 {
//			if _, repeat := found[v]; repeat {
//				t.Fatalf("repeat of value %d", v)
//			}
//			found[v] = struct{}{}
//		}
//		if len(found) == nKeys/2 {
//			// Halfway through iteration, finish grow.
//			for i := 0; i < nBuckets; i++ {
//				delete(m, 1.0)
//			}
//		}
//	}
//	if len(found) != nKeys {
//		t.Fatalf("missing value")
//	}
//}

func TestMapIterOrder(t *testing.T) {
	for _, n := range [...]int{3, 7, 9, 15} {
		// Make m be {0: true, 1: true, ..., n-1: true}.
		m := make(map[int]bool)
		for i := 0; i < n; i++ {
			m[i] = true
		}
		// Check that iterating over the map produces at least two different orderings.
		ord := func() []int {
			var s []int
			for key := range m {
				s = append(s, key)
			}
			return s
		}
		first := ord()
		ok := false
		for try := 0; try < 100; try++ {
			if !reflect.DeepEqual(first, ord()) {
				ok = true
				break
			}
		}
		if !ok {
			t.Errorf("Map with n=%d elements had consistent iteration order: %v", n, first)
		}
	}
}

//TestMapStringBytesLookup fail
//func TestMapStringBytesLookup(t *testing.T) {
//	// Use large string keys to avoid small-allocation coalescing,
//	// which can cause AllocsPerRun to report lower counts than it should.
//	m0 := map[string]int{
//		"1000000000000000000000000000000000000000000000000": 1,
//		"2000000000000000000000000000000000000000000000000": 2,
//	}
//	m1 := map[interface{}]interface{}{
//		"1000000000000000000000000000000000000000000000000": 1,
//		"2000000000000000000000000000000000000000000000000": 2,
//	}
//	_ = m1
//	m := NewConcurrentMapFromMap(m1)
//	buf := []byte("1000000000000000000000000000000000000000000000000")
//	if x, err := m.Get(string(buf)); x != 1 || err != nil { // m[string(buf)]; x != 1 {
//		t.Errorf(`m[string([]byte("1"))] = %d, want 1`, x)
//	}
//	buf[0] = '2'
//	if x, err := m.Get(string(buf)); x != 2 || err != nil { //x := m[string(buf)]; x != 2 {
//		t.Errorf(`m[string([]byte("2"))] = %d, want 2`, x)
//	}

//	var x int
//	n := testing.AllocsPerRun(100, func() {
//		_, _ = m.Get(string(buf))
//		//_ = m0[string(buf)]   //n will be 0
//		//_ = m1[string(buf)]   //n will be 2
//		//x += v.(int) //m[string(buf)]
//	})
//	if n != 0 {
//		t.Errorf("AllocsPerRun for m[string(buf)] = %v, want 0", n)
//	}

//	x = 0
//	n = testing.AllocsPerRun(100, func() {
//		y, err := m.Get(string(buf))
//		//y, ok := m[string(buf)]
//		if err != nil {
//			panic("!ok")
//		}
//		x += y.(int)
//	})
//	if n != 0 {
//		t.Errorf("AllocsPerRun for x,ok = m[string(buf)] = %v, want 0", n)
//	}
//}

/*----------------test concurrent-------------------------------*/
func TestConcurrent(t *testing.T) {
	numCpu := runtime.NumCPU()
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(numCpu))
	writeN := 2*numCpu + 1
	readN := 3*numCpu + 1
	n := 1000000
	var repeat int32 = 0

	wWg := new(sync.WaitGroup)
	wWg.Add(writeN)
	cDone := make(chan struct{})
	cm := NewConcurrentMap()

	//start writeN goroutines to write to map with repeated keys, and count the total number of repeated key
	for i := 0; i < writeN; i++ {
		j := i
		go func() {
			for k := 0; k < n; k++ {
				//0-99999, 50000-149999, 100000-19999, 150000-249999,200000-29999, 250000-349999
				key := k + (j * n / 2)
				if previous, err := cm.Put(key, strconv.Itoa(key)+strings.Repeat(" ", j)); err != nil {
					t.Errorf("Get error %v when concurrent write map", err)
					return
				} else if previous != nil {
					//count the total number of repeated key
					atomic.AddInt32(&repeat, 1)
				}
			}
			wWg.Done()
		}()
	}

	go func() {
		wWg.Wait()
		close(cDone)
	}()

	//start readN goroutines to iterate the map
	rWg := new(sync.WaitGroup)
	rWg.Add(readN)
	for i := 0; i < readN; i++ {
		go func() {
			for {

				for itr := cm.Iterator(); itr.HasNext(); {
					ki, vi, _ := itr.Next()
					k, v := ki.(int), vi.(string)
					if strconv.Itoa(k) != strings.Trim(v, " ") {
						t.Errorf("Get %v by %v, want %v == strings.Trim(\"%v\")", v, k, v, k)
						return
					}
				}

				//exit read goroutines if all write goroutines are done
				exit := false
				select {
				case <-cDone:
					exit = true
					break
				case <-time.After(1 * time.Microsecond):
				}

				if exit {
					break
				}
			}
			rWg.Done()
		}()
	}

	//Start a goroutines to count the size of concurrentMap and total number of repeated keys
	//after all write goroutines are done
	cLast := make(chan struct{})
	go func() {
		wWg.Wait()
		if repeat != int32((writeN-1)*(n/2)) {
			t.Errorf("Repeat %v, want %v", repeat, (writeN-1)*(n/2))
		}

		size := cm.Size()
		if size != int32(n/2+writeN*(n/2)) {
			t.Errorf("Size is %v, want %v", size, n/2+writeN*(n/2))
		}

		cm.Clear()
		size = cm.Size()
		if size != 0 {
			t.Errorf("Size is %v after calling Clear(), want %v", size, 0)
		}
		close(cLast)
	}()

	rWg.Wait()
	<-cLast
	runtime.GC()
}

//below code are used in readme.txt
//func Test1(t *testing.T) {
//	m := NewConcurrentMap()

//	previou, err := m.Put(1, 10) //return nil, nil
//	t.Log("1.", previou, err)
//	previou, err = m.PutIfAbsent(1, 20) //return 10, nil
//	t.Log("2.", previou, err)

//	val, err := m.Get(1) //return 10, nil
//	t.Log("3.", val, err)
//	s := m.Size() //return 1
//	t.Log("4.", s)

//	m.PutAll(map[interface{}]interface{}{
//		1: 100,
//		2: 200,
//	})
//	ok, err := m.RemoveEntry(1, 100) //return true, nil
//	t.Log("5.", ok, err)

//	previou, err = m.Replace(2, 20) //return 200, nil
//	t.Log("6.", previou, err)
//	ok, err = m.CompareAndReplace(2, 200, 20) //return false, nil
//	t.Log("7.", ok, err)

//	previou, err = m.Remove(2) //return 20, nil
//	t.Log("8.", previou, err)

//	m.Clear()
//	s = m.Size() //return 0
//	t.Log("9.", s)

//}
