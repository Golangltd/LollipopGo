package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("-----------", err)
		return
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//	JsonParse := NewJsonStruct()
//	v = Conf.Config{}
//	JsonParse.Load("./conf/config.json", &v)
