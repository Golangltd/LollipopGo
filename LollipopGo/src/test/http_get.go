package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var daat *EmailST

type EmailST struct {
	ID        int
	Sender    string
	Name      string
	Type      int
	Time      string
	Content   string
	IsAdd_ons bool // 是否有附件
	IsOpen    bool // 是否打开过
	IsGet     bool // 是否打开过
	//ItemList  map[int]*ItemST
}

func main() {

	daat := &EmailST{
		ID:        1,
		Sender:    "admin",
		Name:      "admin",
		Type:      1,
		Content:   "qwertyuio",
		IsAdd_ons: false,
	}

	u, _ := url.Parse("http://localhost:8892/GolangLtdGM")
	q := u.Query()
	q.Set("Protocol", "11")
	q.Set("Protocol2", "3")
	q.Set("IMsgtype", "1")
	fmt.Println("---daat", daat)
	str, _ := json.Marshal(&daat)
	q.Set("EmailData", string(str))
	u.RawQuery = q.Encode()
	fmt.Printf("%s /n", u.String())
	res, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("%s", result)
}
