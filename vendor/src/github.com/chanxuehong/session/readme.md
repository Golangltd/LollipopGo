Example:
```Go
package main

import (
	"fmt"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/util/random"
)

type Session struct {
	Id   int64
	Name string
	// TODO: other fields
}

var sessionStore = session.New(60*20, 60*60*24)

func main() {
	ss := &Session{
		Id:   10000,
		Name: "name",
	}

	sid := string(random.NewSessionId())

	if err := sessionStore.Add(sid, ss); err != nil {
		fmt.Println(err)
		return
	}

	// since the stored is pointer, so the modification of Session
	// is not necessary to update sessionStore

	ss.Name = "namex"

	v, err := sessionStore.Get(sid)
	if err != nil {
		fmt.Println(err)
		return
	}
	ssx := v.(*Session)
	fmt.Println(ssx.Name)
}
```
