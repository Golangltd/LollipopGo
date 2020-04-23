package money

import (
	"database/sql"
	"database/sql/driver"
)

var _ sql.Scanner = (*NullMoney)(nil)
var _ driver.Value = NullMoney{}

// NullMoney represents an Money that may be null.
// NullMoney implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullMoney struct {
	Money Money
	Valid bool
}

// Scan implements the sql.Scanner interface.
func (m *NullMoney) Scan(value interface{}) error {
	var n sql.NullInt64
	err := n.Scan(value)
	m.Money = Money(n.Int64)
	m.Valid = n.Valid
	return err
}

// Value implements the driver.Valuer interface.
func (m NullMoney) Value() (driver.Value, error) {
	if !m.Valid {
		return nil, nil
	}
	return int64(m.Money), nil
}
