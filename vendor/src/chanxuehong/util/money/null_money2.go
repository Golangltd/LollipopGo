package money

import (
	"database/sql"
	"database/sql/driver"
)

var _ sql.Scanner = (*NullMoney2)(nil)
var _ driver.Value = NullMoney2{}

// NullMoney2 represents an Money2 that may be null.
// NullMoney2 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullMoney2 struct {
	Money Money2
	Valid bool
}

// Scan implements the sql.Scanner interface.
func (m *NullMoney2) Scan(value interface{}) error {
	if value == nil {
		m.Money, m.Valid = 0, false
		return nil
	}
	m.Valid = true
	return m.Money.Scan(value)
}

// Value implements the driver.Valuer interface.
func (m NullMoney2) Value() (driver.Value, error) {
	if !m.Valid {
		return nil, nil
	}
	return m.Money.Value()
}
