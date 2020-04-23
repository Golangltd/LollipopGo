go-concurrentMap
================

go-concurrentMap is a concurrent Map implement, it is ported from java.util.ConcurrentHashMap.

##### Current version: 1.0 Beta

## Quick start

#### Put, Remove, Replace and Clear methods

```go
m := concurrent.NewConcurrentMap()

previou, err := m.Put(1, 10)                   //return nil, nil
previou, err = m.PutIfAbsent(1, 20)            //return 10, nil

val, err := m.Get(1)                           //return 10, nil
s := m.Size()                                  //return 1

m.PutAll(map[interface{}]interface{}{
	1: 100,
	2: 200,
})

ok, err := m.RemoveEntry(1, 100)               //return true, nil

previou, err = m.Replace(2, 20)                //return 200, nil
ok, err = m.CompareAndReplace(2, 200, 20)      //return false, nil

previou, err = m.Remove(2)                     //return 20, nil

m.Clear()
s = m.Size()                                   //return 0

```

#### Safely use composition operation to update the value from multiple threads

```go
/*---- group string by first char using ConcurrentMap ----*/
//sliceAdd function returns a function that appends v into slice
sliceAdd := func(v interface{}) (func(interface{}) interface{}){
    return func(oldVal interface{})(newVal interface{}){
		if oldVal == nil {
			vs :=  make([]string, 0, 1)
			return append(vs, v.(string))
		} else {
			return append(oldVal.([]string), v.(string))
		}
	}
}

m := concurrent.NewConcurrentMap()
//group by first char of str
group := func(str string) {
	m.Update(string(str[0]), sliceAdd(str))
}

go group("stone")
go group("jack")
go group("jackson")

/*m will include the below key-value pairs, but please note sequence may be different:
{
  s:[stone],
  j:[jack jackson],
}
*/
```

#### Use Hashable interface to customize hash code and equals logic and support reference type and pointer type

```go
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

m := concurrent.NewConcurrentMap()
previou, err := m.Put(&user, 10)                   //return nil, nil
val, err := m.Get(&user)                           //return 10, nil
```

#### Iterator and get key-value slice

```go
//iterate ConcurrentMap
for itr := m.Iterator();itr.HasNext(); {
	k, v, _ := itr.Next()
}

//only user Next method to iterate ConcurrentMap
for itr := m.Iterator();; {
	k, v, ok := itr.Next()
	if !ok {
		break
	}
}

//ToSlice
for _, entry := range m.ToSlice(){
	k, v := entry.Key(), entry.Value()
}
```

#### More factory functions

```go
//new concurrentMap with specified initial capacity
m = concurrent.NewConcurrentMap(32)

//new concurrentMap with specified initial capacity and load factor
m = concurrent.NewConcurrentMap(32, 0.75)

//new concurrentMap with specified initial capacity, load factor and concurrent level
m = concurrent.NewConcurrentMap(32, 0.75, 16)

//new concurrentMap with the same mappings as the given map
m = concurrent.NewConcurrentMapFromMap(map[interface{}]interface{}{
		"x":                      "x1val",
		"xx":                     "x2val",
	})
	
```

## Doc

[Go Doc at godoc.org](https://godoc.org/github.com/fanliao/go-concurrentMap)

## Limitations

* Do not support the below types as key:

   - pointer
   - slice (do not support == operator)
   - map (do not support == operator)
   - channel 
   - function (do not support == operator)
   - struct that includes field which type is above-mentioned or interface
   - array which element type is above-mentioned or interface

   Do not support pointer because the memory address of pointer may be changed after GC, so cannot get a invariant value as hash code for pointer type. Please refer to [when  in next releases of go compacting GC move pointers, does map on poiner types will work ?](https://groups.google.com/forum/#!topic/golang-nuts/AFEf6VM-qrY)

## Performance

Below are the CPU, OS and parameters of benchmark testing: 
 
Xeon E3-1230V3 3.30GHZ, Win7 64 OS

Use 8 procs and 9 goroutines, every goroutines will put or get 100,000 key-value pairs.

I used a thread safe implement that uses the RWMutex to compare the performance,. The below are the test results:

<table>
	<tr>
		<th></th>
		<th>Use RWMutex</th>
		<th>ConcurrentMap</th>
	</tr>
	<tr>
		<td>Put</td>
		<td>480.000 ms/op</td>
		<td>130.207 ms/op</td>
	</tr>
		<td>Get</td>
		<td>45.643 ms/op</td>
		<td>69.464 ms/op</td>
	<tr>
		<td>Put and Get</td>
		<td>729.534 ms/op</td>
		<td>166.610 ms/op</td>
	</tr>
</table>

Note the performance of LockMap's Get operation is better than concurrentMap, the reason is that RWMutex supports parallel read. But if multiple threads put and get at same time, ConcurrentMap will be better than LockMap.

## License

go-concurrentMap is licensed under the MIT Licence, (http://www.apache.org/licenses/LICENSE-2.0.html).
