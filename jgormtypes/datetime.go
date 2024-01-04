package jgormtypes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"time"
)

const ctLayout = "2006-01-02T15:04:05"

var _ JTypeBase = (*JDateTime)(nil)

// JDateTime 时间日期
type JDateTime struct {
	time.Time
	Valid bool
}

func (ct JDateTime) String() string {
	if ct.IsValid() {
		return ct.Time.Format(ctLayout)
	}
	return ""
}

func (dt JDateTime) IsValid() bool {
	return dt.Valid && !dt.Time.IsZero()
}
func (dt JDateTime) IsZero() bool {
	return !dt.IsValid()
}

func (dt JDateTime) MarshalText() ([]byte, error) {
	return []byte(dt.String()), nil
}

func (dt *JDateTime) UnmarshalText(data []byte) error {
	t1, err := parseDateTimeInner(data)
	*dt = JDateTime{
		Time:  t1,
		Valid: err == nil && !t1.IsZero(),
	}
	return err
}

func (ct *JDateTime) UnmarshalJSON(data []byte) (err error) {
	t1, err := parseDateTimeInner(data)
	*ct = JDateTime{
		Time:  t1,
		Valid: err == nil && !t1.IsZero(),
	}
	return
}

func (ct JDateTime) MarshalJSON() ([]byte, error) {
	s1 := ct.String()
	if s1 == "" {
		return []byte("null"), nil
	}
	return bytes.Join([][]byte{{'"'}, []byte(s1), {'"'}}, nil), nil
}

func (n *JDateTime) Scan(value any) error {
	var a sql.NullTime
	err := a.Scan(value)
	if err == nil {
		*n = JDateTime{Time: a.Time, Valid: a.Valid}
	}
	return err
}

func (n JDateTime) Value() (driver.Value, error) {
	if n.IsValid() {
		return n.Time, nil
	}
	return nil, nil
}
