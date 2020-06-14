package jsonutils

import "github.com/pkg/errors"
import "github.com/globalsign/mgo/bson"

//一个动态的json对象
//注意：json在Unmarshal到interface{}时，会把JsonNumber转成float64，除非使用UseNumber
//因此这里仅提供float64接口，其他数据类型外部转换
//如果json的是{type: 1, data: {}}这种格式，需要通过type解析具体的data，则推荐使用json.RawMessage来解析

type JsonObject map[string]interface{}

var (
	TypeError  = errors.New("type convert error")
	KeyError   = errors.New("key not exist")
	IndexError = errors.New("index not exist")
)

func (jm JsonObject) HasKey(key string) bool {
	if _, ok := jm[key]; ok {
		return true
	}
	return false
}

func (jm JsonObject) HasNotNilKey(key string) bool {
	if tmp, ok := jm[key]; ok {
		if tmp != nil {
			return true
		}
	}
	return false
}

func (jm JsonObject) GetObjectId() (bson.ObjectId, error) {
	if tmp, ok := jm["_id"]; ok {
		if id, ok := tmp.(bson.ObjectId); ok {
			return id, nil
		}
	}
	return "", TypeError
}

func (jm JsonObject) GetFloat64(key string) (float64, error) {
	var (
		tmp  interface{}
		resp float64
		ok   bool
	)
	if tmp, ok = jm[key]; ok {
		if resp, ok = tmp.(float64); ok {
			return resp, nil
		}
		return 0, TypeError
	}
	return 0, KeyError
}

func (jm JsonObject) GetFloat64Default(key string, defaultValue float64) float64 {
	var (
		tmp  interface{}
		resp float64
		ok   bool
	)
	if tmp, ok = jm[key]; ok {
		if resp, ok = tmp.(float64); ok {
			return resp
		}
	}
	return defaultValue
}

func (jm JsonObject) GetString(key string) (string, error) {
	var (
		tmp  interface{}
		resp string
		ok   bool
	)
	if tmp, ok = jm[key]; ok {
		if resp, ok = tmp.(string); ok {
			return resp, nil
		}
		return "", TypeError
	}
	return "", KeyError
}

func (jm JsonObject) GetStringDefault(key string, defaultValue string) string {
	var (
		tmp  interface{}
		resp string
		ok   bool
	)
	if tmp, ok = jm[key]; ok {
		if resp, ok = tmp.(string); ok {
			return resp
		}
	}
	return defaultValue
}

func (jm JsonObject) GetBool(key string) (bool, error) {
	var (
		tmp  interface{}
		resp bool
		ok   bool
	)
	if tmp, ok = jm[key]; ok {
		if resp, ok = tmp.(bool); ok {
			return resp, nil
		}
		return false, TypeError
	}
	return false, KeyError
}

func (jm JsonObject) GetBoolDefault(key string, defaultValue bool) bool {
	var (
		tmp  interface{}
		resp bool
		ok   bool
	)
	if tmp, ok = jm[key]; ok {
		if resp, ok = tmp.(bool); ok {
			return resp
		}
	}
	return defaultValue
}

func (jm JsonObject) GetJsonArray(key string) (JsonArray, error) {
	var (
		tmp  interface{}
		resp []interface{}
		ok   bool
	)
	if tmp, ok = jm[key]; ok {
		if resp, ok = tmp.([]interface{}); ok {
			return resp, nil
		}
		return nil, TypeError
	}
	return nil, KeyError
}

func (jm JsonObject) GetJsonObject(key string) (JsonObject, error) {
	var (
		tmp  interface{}
		resp map[string]interface{}
		ok   bool
	)
	if tmp, ok = jm[key]; ok {
		if resp, ok = tmp.(map[string]interface{}); ok {
			return resp, nil
		}
		return nil, TypeError
	}
	return nil, KeyError
}
