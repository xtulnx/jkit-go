package jgormtypes

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"time"
)

const cdLayout = "2006-01-02"

var _ JTypeBase = (*JDate)(nil)

type JDate struct {
	time.Time
	Valid bool
}

func NewDate(y int, m time.Month, d int) JDate {
	return JDate{
		Time:  time.Date(y, m, d, 0, 0, 0, 0, time.Local),
		Valid: true,
	}
}

func (ct JDate) String() string {
	if ct.IsValid() {
		return ct.Time.Format(cdLayout)
	}
	return ""
}

func (dt JDate) IsValid() bool {
	return dt.Valid && !dt.Time.IsZero()
}

func (dt JDate) IsZero() bool {
	return !dt.IsValid()
}

func (dt JDate) MarshalText() ([]byte, error) {
	return []byte(dt.String()), nil
}

func (dt *JDate) UnmarshalText(data []byte) error {
	t1, err := parseDateTimeInner(data)
	y, m, d := t1.Date()
	*dt = JDate{
		Time:  time.Date(y, m, d, 0, 0, 0, 0, t1.Location()),
		Valid: err == nil && !t1.IsZero(),
	}
	return err
}

func (ct *JDate) UnmarshalJSON(b []byte) (err error) {
	t1, err := parseDateTimeInner(b)
	y, m, d := t1.Date()
	*ct = JDate{
		Time:  time.Date(y, m, d, 0, 0, 0, 0, t1.Location()),
		Valid: err == nil && !t1.IsZero(),
	}
	return
}

func (ct JDate) MarshalJSON() ([]byte, error) {
	s1 := ct.String()
	if s1 == "" {
		return []byte("null"), nil
	}
	return bytes.Join([][]byte{{'"'}, []byte(s1), {'"'}}, nil), nil
}

func (n *JDate) Scan(value any) error {
	switch v := value.(type) {
	case []byte:
		if e1 := n.UnmarshalText(v); e1 == nil {
			return nil
		}
	case string:
		if e1 := n.UnmarshalText([]byte(v)); e1 == nil {
			return nil
		}
	}
	var a sql.NullTime
	err := a.Scan(value)
	if err == nil {
		*n = JDate{Time: a.Time, Valid: a.Valid}
	}
	return err
}

func (n JDate) Value() (driver.Value, error) {
	if n.IsValid() {
		return n.Time, nil
	}
	return nil, nil
}

//////////////////////////////////////////////////////////////////

type JDateSlice []JDate

func (s JDateSlice) Len() int {
	return len(s)
}

func (s JDateSlice) Less(i, j int) bool {
	return s[i].Time.Before(s[j].Time)
}
func (s JDateSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
