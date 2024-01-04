package jgormtypes

import (
	"database/sql/driver"
	"github.com/xtulnx/jkit-go/jtime"
	"strings"
	"time"
)

type JTypeBase interface {
	String() string
	IsValid() bool
	IsZero() bool
	UnmarshalText(data []byte) error
	MarshalText() ([]byte, error)
	UnmarshalJSON(data []byte) (err error)
	MarshalJSON() ([]byte, error)
	Scan(value any) error
	Value() (driver.Value, error)
}

func parseDateTimeInner(data []byte) (time.Time, error) {
	if len(data) == 0 {
		return time.Time{}, nil
	}
	s := strings.Trim(string(data), "\"")
	if s == "null" || s == "" {
		return time.Time{}, nil
	}
	if f := dateTimeParser; f != nil {
		return f(s)
	}
	return jtime.StringToDate(s)
}

var dateTimeParser func(string) (time.Time, error) = nil

// 自定义时间日期解析器
func SetDateParser(fn func(string) (time.Time, error)) {
	dateTimeParser = fn
}
