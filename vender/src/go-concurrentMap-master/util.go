package concurrent

import (
	"fmt"
	"hash/fnv"
	"io"
	"math/rand"
	"reflect"
	"sync/atomic"
	"unsafe"
)

const (
	intSize   = unsafe.Sizeof(1)
	ptrSize   = unsafe.Sizeof((*int)(nil))
	bigEndian = false
)

var (
	hasherT           = reflect.TypeOf((*Hashable)(nil)).Elem()
	defaultEqualsfunc func(k1 interface{}, k2 interface{}) bool
	hasherEng         *hashEnginer
	boolEng           *hashEnginer
	intEng            *hashEnginer
	int8Eng           *hashEnginer
	int16Eng          *hashEnginer
	int32Eng          *hashEnginer
	int64Eng          *hashEnginer
	uintEng           *hashEnginer
	uint8Eng          *hashEnginer
	uint16Eng         *hashEnginer
	uint32Eng         *hashEnginer
	uint64Eng         *hashEnginer
	uintptrEng        *hashEnginer
	float32Eng        *hashEnginer
	float64Eng        *hashEnginer
	complex64Eng      *hashEnginer
	complex128Eng     *hashEnginer
	stringEng         *hashEnginer
	engM              map[reflect.Kind]*hashEnginer
)

func init() {
	hasherEng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			w.Write(k.(Hashable).HashBytes())
		},
	}
	boolEng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(bool)
			w.Write((*((*[1]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	intEng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(int)
			w.Write((*((*[intSize]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	int8Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(int8)
			w.Write((*((*[1]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	int16Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(int16)
			w.Write((*((*[2]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	int32Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(int32)
			w.Write((*((*[4]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	int64Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(int64)
			w.Write((*((*[8]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	uintEng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(uint)
			w.Write((*((*[intSize]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	uint8Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(uint8)
			w.Write((*((*[1]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	uint16Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(uint16)
			w.Write((*((*[2]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	uint32Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(uint32)
			w.Write((*((*[4]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	uint64Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(uint64)
			w.Write((*((*[8]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	uintptrEng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(uintptr)
			w.Write((*((*[intSize]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	float32Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(float32)
			//Nan != Nan, so use a rand number to generate hash code
			if k1 != k1 {
				k1 = rand.Float32()
			}
			w.Write((*((*[4]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	float64Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(float64)
			//Nan != Nan, so use a rand number to generate hash code
			if k1 != k1 {
				k1 = rand.Float64()
			}
			w.Write((*((*[8]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	complex64Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(complex64)
			w.Write((*((*[8]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	complex128Eng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(complex128)
			w.Write((*((*[128]byte)(unsafe.Pointer(&k1))))[:])
		},
	}
	stringEng = &hashEnginer{
		putFunc: func(w io.Writer, k interface{}) {
			k1 := k.(string)
			w.Write([]byte(k1))
		},
	}
	engM = map[reflect.Kind]*hashEnginer{
		reflect.Bool:       boolEng,
		reflect.Int:        intEng,
		reflect.Int8:       int8Eng,
		reflect.Int16:      int16Eng,
		reflect.Int32:      int32Eng,
		reflect.Int64:      int64Eng,
		reflect.Uint:       uintEng,
		reflect.Uint8:      uint8Eng,
		reflect.Uint16:     uint16Eng,
		reflect.Uint32:     uint32Eng,
		reflect.Uint64:     uint64Eng,
		reflect.Uintptr:    uintptrEng,
		reflect.Float32:    float32Eng,
		reflect.Float64:    float64Eng,
		reflect.Complex64:  complex64Eng,
		reflect.Complex128: complex128Eng,
		reflect.String:     stringEng,
	}
}

func hashKey(key interface{}, m *ConcurrentMap, isRead bool) (hashCode uint32, err error) {
	h := fnv.New32a()

	switch v := key.(type) {
	case bool:
		h.Write((*((*[1]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case int:
		h.Write((*((*[intSize]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case int8:
		h.Write((*((*[1]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case int16:
		h.Write((*((*[2]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case int32:
		h.Write((*((*[4]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case int64:
		h.Write((*((*[8]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case uint:
		h.Write((*((*[intSize]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case uint8:
		h.Write((*((*[1]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case uint16:
		h.Write((*((*[2]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case uint32:
		h.Write((*((*[4]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case uint64:
		h.Write((*((*[8]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case uintptr:
		h.Write((*((*[intSize]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case float32:
		//Nan != Nan, so use a rand number to generate hash code
		if v != v {
			v = rand.Float32()
		}
		h.Write((*((*[4]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case float64:
		//Nan != Nan, so use a rand number to generate hash code
		if v != v {
			v = rand.Float64()
		}
		h.Write((*((*[8]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case complex64:
		h.Write((*((*[8]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case complex128:
		h.Write((*((*[128]byte)(unsafe.Pointer(&v))))[:])
		hashCode = h.Sum32()
	case string:
		h.Write([]byte(v))
		hashCode = h.Sum32()
	default:
		//if key is not simple type
		if her, ok := key.(Hashable); ok {
			h.Write(her.HashBytes())
		} else {
			if err = m.parseKey(key); err != nil {
				return
			}
			if isRead {
				eng := (*hashEnginer)(atomic.LoadPointer(&m.eng))
				eng.putFunc(h, key)
			} else {
				eng := (*hashEnginer)(m.eng)
				eng.putFunc(h, key)
			}
			hashCode = h.Sum32()
		}
	}
	return
}

////hash a interface using FNVa
//func hashI(val interface{}) (hashCode uint32) {
//	h := fnv.New32a()
//	switch v := val.(type) {
//	case bool:
//		h.Write((*((*[1]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case int:
//		h.Write((*((*[intSize]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case int8:
//		h.Write((*((*[1]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case int16:
//		h.Write((*((*[2]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case int32:
//		h.Write((*((*[4]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case int64:
//		h.Write((*((*[8]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case uint:
//		h.Write((*((*[1]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case uint8:
//		h.Write((*((*[intSize]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case uint16:
//		h.Write((*((*[2]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case uint32:
//		h.Write((*((*[4]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case uint64:
//		h.Write((*((*[8]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case uintptr:
//		h.Write((*((*[intSize]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case float32:
//		//Nan != Nan, so use a rand number to generate hash code
//		if v != v {
//			v = rand.Float32()
//		}
//		h.Write((*((*[4]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case float64:
//		//Nan != Nan, so use a rand number to generate hash code
//		if v != v {
//			v = rand.Float64()
//		}
//		h.Write((*((*[8]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case complex64:
//		h.Write((*((*[8]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case complex128:
//		h.Write((*((*[128]byte)(unsafe.Pointer(&v))))[:])
//		hashCode = h.Sum32()
//	case string:
//		h.Write([]byte(v))
//		hashCode = h.Sum32()
//	default:
//		//some types can be used as key, we can use equals to test
//		_ = val == val

//		//support array, struct, channel, interface, pointer
//		//don't support slice, function, map
//		rv := reflect.ValueOf(val)
//		switch rv.Kind() {
//		case reflect.Ptr:
//			//ei.word stores the memory address of value that v points to, we use address to generate hash code
//			ei := (*emptyInterface)(unsafe.Pointer(&val))
//			hashCode = hashI(uintptr(ei.word))
//		case reflect.Interface:
//			//for interface, we use contained value to generate the hash code
//			hashCode = hashI(rv.Elem())
//		default:
//			//for array, struct and chan, will get byte array to calculate the hash code
//			hashMem(rv, h)
//			hashCode = h.Sum32()
//			fmt.Println("array, struct or chan", rv.Interface(), hashCode, reflect.ValueOf(rv).Type().Size())
//		}
//	}
//	return
//}

////hashMem writes byte array of underlying value to hash function
//func hashMem(i interface{}, hashFunc hash.Hash32) {
//	fmt.Println("hashMem")
//	size := reflect.ValueOf(i).Type().Size()
//	ei := (*emptyInterface)(unsafe.Pointer(&i))

//	//if size of underlying value is greater than pointer size, ei.word will store the pointer that point to underlying value
//	//else ei.word will store underlying value
//	if size > ptrSize {
//		addr := ei.word
//		hashPtrData(unsafe.Pointer(uintptr(addr)), size, hashFunc)
//	} else {
//		data := ei.word
//		fmt.Println("hashData", uintptr(data), size, ptrSize)
//		hashData(uintptr(data), size, hashFunc)
//	}
//	return
//}

//func hashPtrData(basePtr unsafe.Pointer, size uintptr, hashFunc hash.Hash32) {
//	offset := uintptr(0)
//	for {
//		/* cannot store unsafe.Pointer in an uintptr according to https://groups.google.com/forum/#!topic/golang-dev/bfMdPAQigfM
//		 * but the expression
//		 *     unsafe.Pointer(uintptr(basePtr) + offset)
//		 * is safe under Go 1.3
//		 */
//		//d := uintptr(basePtr) + offset
//		//ptr := unsafe.Pointer(d)
//		ptr := unsafe.Pointer(uintptr(basePtr) + offset)

//		if size >= 32 {
//			bytes := *(*[32]byte)(ptr)
//			size -= 32
//			offset += 32
//			fmt.Println("hashPtrData", ptr, bytes[:])
//			hashFunc.Write(bytes[:])
//		} else if size >= 16 {
//			bytes := *(*[16]byte)(ptr)
//			size -= 16
//			offset += 16
//			hashFunc.Write(bytes[:])
//		} else if size >= 8 {
//			bytes := *(*[8]byte)(ptr)
//			size -= 8
//			offset += 8
//			hashFunc.Write(bytes[:])
//		} else if size >= 4 {
//			bytes := *(*[4]byte)(ptr)
//			size -= 4
//			offset += 4
//			hashFunc.Write(bytes[:])
//		} else if size >= 2 {
//			bytes := *(*[2]byte)(ptr)
//			size -= 2
//			offset += 2
//			hashFunc.Write(bytes[:])
//		} else if size == 1 {
//			bytes := *(*[1]byte)(ptr)
//			hashFunc.Write(bytes[:])
//			return
//		}
//		if size == 0 {
//			return
//		}
//	}
//}

//func hashData(data uintptr, size uintptr, hashFunc hash.Hash32) {
//	bytes := (*((*[ptrSize]byte)(unsafe.Pointer(&data))))
//	hashFunc.Write(bytes[0:size])
//	return
//}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	k := rv.Type().Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

//// emptyInterface is the header for an interface{} value.
//type emptyInterface struct {
//	typ  uintptr
//	word unsafe.Pointer
//}

type keyInfo struct {
	isHasher bool
	/*-- kind of key type --*/
	kind reflect.Kind
	/*-- index of field if it is a field of struct --*/
	index []int
	/*-- field informations of struct --*/
	fields []*keyInfo
	/*-- element information of array --*/
	elementInfo *keyInfo
	size        int
}

//获取t对应的类型信息，不支持slice, function, map, pointer, interface, channel
func getKeyInfo(t reflect.Type) (ki *keyInfo, err error) {
	return getKeyInfoByParent(t, nil, make([]int, 0, 0))
}

//获取t对应的类型信息，不支持slice, function, map, pointer, interface, channel
//如果parentIdx的长度>0，则表示t是strut中的字段的类型信息, t为字段对应的类型
func getKeyInfoByParent(t reflect.Type, parent *keyInfo, parentIdx []int) (ki *keyInfo, err error) {
	ki = &keyInfo{}
	//判断是否实现了hasher接口
	if t.Implements(hasherT) {
		ki.isHasher = true
		return
	}
	ki.kind = t.Kind()

	if _, ok := engM[ki.kind]; ok {
		//简单类型，不需要再分解元素类型的信息
		ki.index = parentIdx
	} else {
		//some types can be used as key, we can use equals to test
		switch ki.kind {
		case reflect.Chan, reflect.Slice, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface:
			err = NonSupportKey
		case reflect.Struct:
			if parent == nil {
				//parent==nil表示t不是一个嵌套的struct，所以这里需要初始化fields
				parent = ki
				ki.fields = make([]*keyInfo, 0, t.NumField())
			}
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				//skip unexported field,
				if len(f.PkgPath) > 0 {
					continue
				}

				idx := make([]int, len(parentIdx), len(parentIdx)+1)
				copy(idx, parentIdx)
				idx = append(idx, i)
				if fi, e := getKeyInfoByParent(f.Type, parent, idx); e != nil {
					err = e
					return
				} else {
					//fi.index = i
					parent.fields = append(ki.fields, fi)
				}
			}
		case reflect.Array:
			if ki.elementInfo, err = getKeyInfo(t.Elem()); err != nil {
				return
			}
			ki.size = t.Len()
			ki.index = parentIdx
		}
	}

	return

}

func getPutFunc(ki *keyInfo) func(w io.Writer, k interface{}) {
	if ki.isHasher {
		return hasherEng.putFunc
	}

	//Printf("getPutFunc, ki = %v\n", ki)
	if eng, ok := engM[ki.kind]; ok {
		return eng.putFunc
	} else {
		if ki.kind == reflect.Struct {
			//Printf("getPutFunc, ki = %v, other case\n", ki)
			putFunc := func(w io.Writer, k interface{}) {
				rv := reflect.ValueOf(k)
				for _, fieldInfo := range ki.fields {
					//深度遍历每个field，并将其[]byte写入hash函数
					putF := getPutFunc(fieldInfo)
					//Printf("getPutFunc, ki = %#v, fieldInfo = %#v, %#v\n", ki, fieldInfo, rv.Interface())
					//Printf("getPutFunc, value = %v\n", rv.FieldByIndex(fieldInfo.index).Interface())
					putF(w, rv.FieldByIndex(fieldInfo.index).Interface())
				}
			}
			//Printf("getPutFunc, ki=%v, putFunc = %v, other case\n", ki, putFunc)
			return putFunc
		} else if ki.kind == reflect.Array {
			putFunc := func(w io.Writer, k interface{}) {
				rv := reflect.ValueOf(k)
				putF := getPutFunc(ki.elementInfo)
				for i := 0; i < ki.size; i++ {
					//遍历数组元素，并将其[]byte写入hash函数
					putF(w, rv.Index(i).Interface())
				}
			}
			return putFunc
		}
		Printf("getPutFunc, return nil")
	}
	return nil
}

func equals(k1, k2 interface{}) bool {
	if h1, ok := k1.(Hashable); ok {
		return h1.Equals(k2)
	} else {
		return k1 == k2
	}
}

func Printf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		return fmt.Printf(format, a...)
	} else {
		return 0, nil
	}
}

func Println(a ...interface{}) (n int, err error) {
	if Debug {
		return fmt.Println(a...)
	} else {
		return 0, nil
	}
}
