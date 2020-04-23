package json

import (
	"errors"
	"fmt"
	"strconv"
)

type Int64 int64

func (x Int64) MarshalJSON() (data []byte, err error) {
	data = make([]byte, 0, 20+2)
	data = append(data, '"')
	data = strconv.AppendInt(data, int64(x), 10)
	data = append(data, '"')
	return
}

func (x *Int64) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("json: cannot unmarshal empty string into Go value of type Int64")
	}
	if len(data) > 20+2 {
		return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Int64", data)
	}
	if data[0] != '"' {
		n, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Int64", data)
		}
		*x = Int64(n)
		return nil
	}
	maxIndex := len(data) - 1
	if maxIndex < 2 || data[maxIndex] != '"' {
		return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Int64", data)
	}
	n, err := strconv.ParseInt(string(data[1:maxIndex]), 10, 64)
	if err != nil {
		return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Int64", data)
	}
	*x = Int64(n)
	return nil
}
