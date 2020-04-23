package json

import (
	"errors"
	"fmt"
	"strconv"
)

type Uint uint

func (x Uint) MarshalJSON() (data []byte, err error) {
	data = make([]byte, 0, 20+2)
	data = append(data, '"')
	data = strconv.AppendUint(data, uint64(x), 10)
	data = append(data, '"')
	return
}

func (x *Uint) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("json: cannot unmarshal empty string into Go value of type Uint")
	}
	if len(data) > 20+2 {
		return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Uint", data)
	}
	if data[0] != '"' {
		n, err := strconv.ParseUint(string(data), 10, 0)
		if err != nil {
			return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Uint", data)
		}
		*x = Uint(n)
		return nil
	}
	maxIndex := len(data) - 1
	if maxIndex < 2 || data[maxIndex] != '"' {
		return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Uint", data)
	}
	n, err := strconv.ParseUint(string(data[1:maxIndex]), 10, 0)
	if err != nil {
		return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Uint", data)
	}
	*x = Uint(n)
	return nil
}
