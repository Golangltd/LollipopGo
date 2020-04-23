// session implements a simple memory-based session container.
// @link        https://github.com/chanxuehong/session for the canonical source repository
// @license     https://github.com/chanxuehong/session/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package session

import (
	"bytes"
	"testing"
	"time"
)

// Convert Storage to an array, assuming storage types are byte
func convert2array(s *Storage) []byte {
	if len(s.cache) != s.lruList.Len() {
		panic("len(s.cache) != s.lruList.Len()")
	}

	arr := make([]byte, s.lruList.Len())
	for i, e := 0, s.lruList.Front(); e != nil; i, e = i+1, e.Next() {
		payload := e.Value.(*payload)

		if s.cache[payload.Key] != e {
			panic("s.cache[e.Value.(*payload).Key] != e")
		}
		arr[i] = payload.Value.(byte)
	}

	return arr
}

// no expiry
func TestStorageAdd1(t *testing.T) {
	s := New(1, 10)

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("2", byte(2)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("3", byte(3)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("3", byte(3)); err != nil {
		if err != ErrNotStored {
			t.Error("the err must be ErrNotStored")
			return
		}
	} else {
		t.Error("Add duplicate keys should go wrong")
		return
	}
	if err := s.Add("4", byte(4)); err != nil {
		t.Error(err)
		return
	}

	have := convert2array(s)
	want := []byte{4, 3, 2, 1}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}
}

// has expiry
func TestStorageAdd2(t *testing.T) {
	s := New(1, 10)

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("2", byte(2)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("3", byte(3)); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 2)

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("4", byte(4)); err != nil {
		t.Error(err)
		return
	}

	have := convert2array(s)
	want := []byte{4, 1, 3}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}
}

// no expiry
func TestStorageDelete1(t *testing.T) {
	s := New(1, 10)

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("2", byte(2)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("3", byte(3)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Delete("4"); err != nil {
		if err != ErrNotFound {
			t.Error("the err must be ErrNotFound")
			return
		}
	} else {
		t.Error("Delete elements that does not exist should go wrong")
		return
	}
	if err := s.Delete("2"); err != nil {
		t.Error(err)
		return
	}

	have := convert2array(s)
	want := []byte{3, 1}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}
}

// remove elements that expired
func TestStorageDelete2(t *testing.T) {
	s := New(1, 10)

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("2", byte(2)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("3", byte(3)); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 2)

	if err := s.Delete("2"); err != nil {
		if err != ErrNotFound {
			t.Error("the err must be ErrNotFound")
			return
		}
	} else {
		t.Error("Delete elements that expired should go wrong")
		return
	}

	have := convert2array(s)
	want := []byte{3, 1}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}
}

// no expiry
func TestStorageGet1(t *testing.T) {
	s := New(1, 10)

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("2", byte(2)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("3", byte(3)); err != nil {
		t.Error(err)
		return
	}
	if _, err := s.Get("4"); err != nil {
		if err != ErrNotFound {
			t.Error("the err must be ErrNotFound")
			return
		}
	} else {
		t.Error("Get elements that does not exist should go wrong")
		return
	}
	if v, err := s.Get("1"); err != nil {
		t.Error(err)
		return
	} else {
		if n := v.(byte); n != 1 {
			t.Error(`Get("1"), have:`, n, "want:", 1)
			return
		}
	}
	if v, err := s.Get("2"); err != nil {
		t.Error(err)
		return
	} else {
		if n := v.(byte); n != 2 {
			t.Error(`Get("2"), have:`, n, "want:", 2)
			return
		}
	}

	have := convert2array(s)
	want := []byte{2, 1, 3}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}
}

// get elements that expired
func TestStorageGet2(t *testing.T) {
	s := New(1, 10)

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("2", byte(2)); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 2)

	if err := s.Add("3", byte(3)); err != nil {
		t.Error(err)
		return
	}

	if _, err := s.Get("2"); err != nil {
		if err != ErrNotFound {
			t.Error("the err must be ErrNotFound")
			return
		}
	} else {
		t.Error("Get elements that expired should go wrong")
		return
	}
	if v, err := s.Get("3"); err != nil {
		t.Error(err)
		return
	} else {
		if n := v.(byte); n != 3 {
			t.Error(`Get("3"), have:`, n, "want:", 3)
			return
		}
	}

	have := convert2array(s)
	want := []byte{3}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}
}

// no expiry
func TestStorageSet1(t *testing.T) {
	s := New(1, 10)

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("2", byte(2)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("3", byte(3)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Set("1", byte(11)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Set("4", byte(4)); err != nil {
		t.Error(err)
		return
	}

	have := convert2array(s)
	want := []byte{4, 11, 3, 2}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}
}

// set elements that expired
func TestStorageSet2(t *testing.T) {
	s := New(1, 10)

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("2", byte(2)); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 2)

	if err := s.Set("1", byte(11)); err != nil {
		t.Error(err)
		return
	}

	have := convert2array(s)
	want := []byte{11, 2}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}
}

func TestStorageGC(t *testing.T) {
	s := New(4, 5)

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("2", byte(2)); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 3)

	// no expiry, no gc() run

	if err := s.Add("3", byte(3)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("4", byte(4)); err != nil {
		t.Error(err)
		return
	}

	have := convert2array(s)
	want := []byte{4, 3, 2, 1}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}

	time.Sleep(time.Second * 3)

	// elements with key "1", "2" has expired, gc() has run

	have = convert2array(s)
	want = []byte{4, 3}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}
}

func TestStorageSetGCInterval(t *testing.T) {
	s := New(4, 100)

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("2", byte(2)); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 5)

	s.SetGCInterval(5) // elements with key "1", "2" has expired, and call gc()

	have := convert2array(s)
	want := []byte{}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}

	// see TestStorageGC

	if err := s.Add("1", byte(1)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("2", byte(2)); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Second * 3)

	if err := s.Add("3", byte(3)); err != nil {
		t.Error(err)
		return
	}
	if err := s.Add("4", byte(4)); err != nil {
		t.Error(err)
		return
	}

	have = convert2array(s)
	want = []byte{4, 3, 2, 1}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}

	time.Sleep(time.Second * 3)

	have = convert2array(s)
	want = []byte{4, 3}
	if !bytes.Equal(have, want) {
		t.Error("have:", have, "want:", want)
		return
	}
}
