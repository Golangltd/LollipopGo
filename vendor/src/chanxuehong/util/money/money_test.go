package money

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

// TestMoneyMarshalJSON also test String and MarshalText
func TestMoneyMarshalJSON(t *testing.T) {
	tests := []struct {
		src Money
		dst string
	}{
		{1, `"0.01"`},
		{11, `"0.11"`},
		{100, `"1"`},
		{110, `"1.10"`},
		{111, `"1.11"`},
		{1111, `"11.11"`},
		{0, `"0"`},
		{-1, `"-0.01"`},
		{-11, `"-0.11"`},
		{-100, `"-1"`},
		{-110, `"-1.10"`},
		{-111, `"-1.11"`},
		{-1111, `"-11.11"`},
	}
	var (
		text []byte
		err  error
	)
	for _, pair := range tests {
		text, err = json.Marshal(pair.src)
		if err != nil {
			t.Errorf("json.Marshal Money %d failed: %s\r\n", int64(pair.src), err.Error())
			continue
		}
		if string(text) != pair.dst {
			t.Errorf("json.Marshal Money %d failed, have %s, want %s\r\n", int64(pair.src), text, pair.dst)
			continue
		}
	}
}

func TestMoneyMarshalXML(t *testing.T) {
	tests := []struct {
		src Money
		dst string
	}{
		{1, `<Money>0.01</Money>`},
		{11, `<Money>0.11</Money>`},
		{100, `<Money>1</Money>`},
		{110, `<Money>1.10</Money>`},
		{111, `<Money>1.11</Money>`},
		{1111, `<Money>11.11</Money>`},
		{0, `<Money>0</Money>`},
		{-1, `<Money>-0.01</Money>`},
		{-11, `<Money>-0.11</Money>`},
		{-100, `<Money>-1</Money>`},
		{-110, `<Money>-1.10</Money>`},
		{-111, `<Money>-1.11</Money>`},
		{-1111, `<Money>-11.11</Money>`},
	}
	var (
		text []byte
		err  error
	)
	for _, pair := range tests {
		text, err = xml.Marshal(pair.src) // The name for the XML elements is the name of the marshalled type
		if err != nil {
			t.Errorf("xml.Marshal Money %d failed: %s\r\n", int64(pair.src), err.Error())
			continue
		}
		if string(text) != pair.dst {
			t.Errorf("xml.Marshal Money %d failed, have %s, want %s\r\n", int64(pair.src), text, pair.dst)
			continue
		}
	}
}

// TestMoneyUnmarshalJSON also test UnmarshalText
func TestMoneyUnmarshalJSON(t *testing.T) {
	var (
		money Money
		err   error
	)

	for _, str := range []string{`"1x"`, `"0.11x"`, `"0.111"`} {
		if err = json.Unmarshal([]byte(str), &money); err == nil {
			t.Errorf("json.Unmarshal %s to Money should failed, but not\r\n", str)
			continue
		}
	}

	tests := []struct {
		dst Money
		src string
	}{
		{1, `"0.01"`},
		{1, `".01"`},
		{11, `"0.11"`},
		{100, `"1"`},
		{100, `"1.00"`},
		{110, `"1.10"`},
		{111, `"1.11"`},
		{1111, `"11.11"`},
		{0, `"0"`},
		{0, `".0"`},
		{0, `".00"`},
		{0, `"0.0"`},
		{-1, `"-.01"`},
		{-1, `"-0.01"`},
		{-11, `"-0.11"`},
		{-100, `"-1"`},
		{-100, `"-1.00"`},
		{-110, `"-1.10"`},
		{-111, `"-1.11"`},
		{-1111, `"-11.11"`},
	}
	for _, pair := range tests {
		if err = json.Unmarshal([]byte(pair.src), &money); err != nil {
			t.Errorf("json.Unmarshal %s to Money failed: %s\r\n", pair.src, err.Error())
			continue
		}
		if money != pair.dst {
			t.Errorf("json.Unmarshal %s to Money failed, have %d, want %d\r\n", pair.src, money, pair.dst)
			continue
		}
	}
}

func TestMoneyUnmarshalXML(t *testing.T) {
	var (
		money Money
		err   error
	)

	for _, str := range []string{`<Money>1x</Money>`, `<Money>0.01x</Money>`, `<Money>0.011</Money>`} {
		if err = xml.Unmarshal([]byte(str), &money); err == nil {
			t.Errorf("xml.Unmarshal %s to Money should failed, but not\r\n", str)
			continue
		}
	}

	tests := []struct {
		dst Money
		src string
	}{
		{1, `<Money>0.01</Money>`},
		{1, `<Money>.01</Money>`},
		{11, `<Money>0.11</Money>`},
		{100, `<Money>1</Money>`},
		{100, `<Money>1.00</Money>`},
		{110, `<Money>1.10</Money>`},
		{111, `<Money>1.11</Money>`},
		{1111, `<Money>11.11</Money>`},
		{0, `<Money>0</Money>`},
		{0, `<Money>.0</Money>`},
		{0, `<Money>.00</Money>`},
		{0, `<Money>0.0</Money>`},
		{-1, `<Money>-.01</Money>`},
		{-1, `<Money>-0.01</Money>`},
		{-11, `<Money>-0.11</Money>`},
		{-100, `<Money>-1</Money>`},
		{-100, `<Money>-1.00</Money>`},
		{-110, `<Money>-1.10</Money>`},
		{-111, `<Money>-1.11</Money>`},
		{-1111, `<Money>-11.11</Money>`},
	}
	for _, pair := range tests {
		if err = xml.Unmarshal([]byte(pair.src), &money); err != nil {
			t.Errorf("xml.Unmarshal %s to Money failed: %s\r\n", pair.src, err.Error())
			continue
		}
		if money != pair.dst {
			t.Errorf("xml.Unmarshal %s to Money failed, have %d, want %d\r\n", pair.src, money, pair.dst)
			continue
		}
	}
}
