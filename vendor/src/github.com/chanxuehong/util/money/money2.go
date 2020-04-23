package money

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
)

// Money2 表示金钱, 单位为分.
//
// Money2 是 Money 的扩展, 支持数据库 decimal(x.2) 类型的直接存取.
type Money2 int64

var (
	_ encoding.TextMarshaler   = Money2(0)
	_ encoding.TextUnmarshaler = (*Money2)(nil)
)

var (
	_ json.Marshaler   = Money2(0)
	_ json.Unmarshaler = (*Money2)(nil)
)

var (
	_ xml.Marshaler   = Money2(0)
	_ xml.Unmarshaler = (*Money2)(nil)
)

var (
	_ driver.Valuer = Money2(0)
	_ sql.Scanner   = (*Money2)(nil)
)

// Value 实现了 driver.Valuer 接口, 将 Money2 编码成 decimal(.2) 格式.
func (m Money2) Value() (driver.Value, error) {
	return m.MarshalText()
}

// Scan 实现了 sql.Scanner 接口, 将 decimal(.2) 字段解码到 Money2 中.
func (m *Money2) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return m.UnmarshalText(v)
	case string:
		return m.UnmarshalText([]byte(v))
	case nil:
		return errors.New("unsupported Scan, storing nil driver.Value into type Money2")
	default:
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type Money2", value)
	}
}

// Text 将 Money2 编码成 xxxx.yz 这样以 '元' 为单位的字符串.
func (m Money2) Text() string {
	return Money(m).Text()
}

// MarshalText 将 Money2 编码成 xxxx.yz 这样以 '元' 为单位的字符串.
func (m Money2) MarshalText() (text []byte, err error) {
	return Money(m).MarshalText()
}

// MarshalJSON 将 Money2 编码成 "xxxx.yz" 这样以 '元' 为单位的字符串.
func (m Money2) MarshalJSON() ([]byte, error) {
	return Money(m).MarshalJSON()
}

// MarshalXML 将 Money2 编码成 xxxx.yz 这样以 '元' 为单位的字符串.
func (m Money2) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	return Money(m).MarshalXML(e, start)
}

// UnmarshalText 将 xxxx.yz 这样以 '元' 为单位的字符串解码到 Money2 中.
func (m *Money2) UnmarshalText(text []byte) (err error) {
	return ((*Money)(m)).UnmarshalText(text)
}

// UnmarshalTextString 将 xxxx.yz 这样以 '元' 为单位的字符串解码到 Money2 中.
func (m *Money2) UnmarshalTextString(text string) (err error) {
	return ((*Money)(m)).UnmarshalTextString(text)
}

// UnmarshalJSON 将 "xxxx.yz" 这样以 '元' 为单位的字符串解码到 Money2 中.
func (m *Money2) UnmarshalJSON(data []byte) (err error) {
	return ((*Money)(m)).UnmarshalJSON(data)
}

// UnmarshalXML 将 xxxx.yz 这样以 '元' 为单位的字符串解码到 Money2 中.
func (m *Money2) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	return ((*Money)(m)).UnmarshalXML(d, start)
}
