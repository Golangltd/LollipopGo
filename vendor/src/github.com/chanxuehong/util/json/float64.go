package json

import (
	"errors"
	"fmt"
	"strconv"
)

type Float64 float64

func (x Float64) MarshalJSON() (data []byte, err error) {
	data = make([]byte, 0, 24+2)
	data = append(data, '"')
	data = strconv.AppendFloat(data, float64(x), 'g', -1, 64)
	data = append(data, '"')
	return
}

func (x *Float64) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("json: cannot unmarshal empty string into Go value of type Float64")
	}
	if data[0] != '"' {
		n, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Float64", data)
		}
		*x = Float64(n)
		return nil
	}
	maxIndex := len(data) - 1
	if maxIndex < 2 || data[maxIndex] != '"' {
		return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Float64", data)
	}
	n, err := strconv.ParseFloat(string(data[1:maxIndex]), 64)
	if err != nil {
		return fmt.Errorf("json: cannot unmarshal string %s into Go value of type Float64", data)
	}
	*x = Float64(n)
	return nil
}
