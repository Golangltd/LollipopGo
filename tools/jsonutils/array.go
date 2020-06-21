package jsonutils

type JsonArray []interface{}

func (ja JsonArray) ToNumberArray() ([]float64, error) {
	if ja == nil {
		return nil, TypeError
	}
	resp := make([]float64, len(ja))
	var (
		v   float64
		err error
	)
	for i := 0; i < len(ja); i++ {
		v, err = ja.GetFloat64ByIndex(i)
		if err != nil {
			return nil, err
		}
		resp = append(resp, v)
	}
	return resp, nil
}

func (ja JsonArray) ToStringArray() ([]string, error) {
	if ja == nil {
		return nil, TypeError
	}
	resp := make([]string, len(ja))
	var (
		v   string
		err error
	)
	for i := 0; i < len(ja); i++ {
		v, err = ja.GetStringByIndex(i)
		if err != nil {
			return nil, err
		}
		resp = append(resp, v)
	}
	return resp, nil
}

func (ja JsonArray) ToBoolArray() ([]bool, error) {
	if ja == nil {
		return nil, TypeError
	}
	resp := make([]bool, len(ja))
	var (
		v   bool
		err error
	)
	for i := 0; i < len(ja); i++ {
		v, err = ja.GetBoolByIndex(i)
		if err != nil {
			return nil, err
		}
		resp = append(resp, v)
	}
	return resp, nil
}

func (ja JsonArray) ToObjectArray() ([]JsonObject, error) {
	if ja == nil {
		return nil, TypeError
	}
	resp := make([]JsonObject, len(ja))
	var (
		v   JsonObject
		err error
	)
	for i := 0; i < len(ja); i++ {
		v, err = ja.GetObjectByIndex(i)
		if err != nil {
			return nil, err
		}
		resp = append(resp, v)
	}
	return resp, nil
}

func (ja JsonArray) ToArrayOfArray() ([]JsonArray, error) {
	if ja == nil {
		return nil, TypeError
	}
	resp := make([]JsonArray, len(ja))
	var (
		v   JsonArray
		err error
	)
	for i := 0; i < len(ja); i++ {
		v, err = ja.GetArrayByIndex(i)
		if err != nil {
			return nil, err
		}
		resp = append(resp, v)
	}
	return resp, nil
}

func (ja JsonArray) GetFloat64ByIndex(index int) (float64, error) {
	var (
		tmp  interface{}
		resp float64
		ok   bool
	)
	if index < 0 || ja == nil || index >= len(ja) {
		return 0, IndexError
	}
	tmp = ja[index]
	if resp, ok = tmp.(float64); ok {
		return resp, nil
	}
	return 0, TypeError
}

func (ja JsonArray) GetStringByIndex(index int) (string, error) {
	var (
		tmp  interface{}
		resp string
		ok   bool
	)
	if index < 0 || ja == nil || index >= len(ja) {
		return "", IndexError
	}
	tmp = ja[index]
	if resp, ok = tmp.(string); ok {
		return resp, nil
	}
	return "", TypeError
}

func (ja JsonArray) GetBoolByIndex(index int) (bool, error) {
	var (
		tmp  interface{}
		resp bool
		ok   bool
	)
	if index < 0 || ja == nil || index >= len(ja) {
		return false, IndexError
	}
	tmp = ja[index]
	if resp, ok = tmp.(bool); ok {
		return resp, nil
	}
	return false, TypeError
}

func (ja JsonArray) GetObjectByIndex(index int) (JsonObject, error) {
	var (
		tmp  interface{}
		resp map[string]interface{}
		ok   bool
	)
	if index < 0 || ja == nil || index >= len(ja) {
		return nil, IndexError
	}
	tmp = ja[index]
	if resp, ok = tmp.(map[string]interface{}); ok {
		return resp, nil
	}
	return nil, TypeError
}

func (ja JsonArray) GetArrayByIndex(index int) (JsonArray, error) {
	var (
		tmp  interface{}
		resp []interface{}
		ok   bool
	)
	if index < 0 || ja == nil || index >= len(ja) {
		return nil, IndexError
	}
	tmp = ja[index]
	if resp, ok = tmp.([]interface{}); ok {
		return resp, nil
	}
	return nil, TypeError
}
