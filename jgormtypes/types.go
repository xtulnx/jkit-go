package jgormtypes

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/xtulnx/jkit-go/jtime"
)

type JTypeBase interface {
	String() string                        // 转换为字符串
	IsValid() bool                         // 是否有效值
	IsZero() bool                          // 是否零值（即无效）
	UnmarshalText(data []byte) error       // 从文本解析
	MarshalText() ([]byte, error)          // 转换为文本
	UnmarshalJSON(data []byte) (err error) // 从JSON解析
	MarshalJSON() ([]byte, error)          // 转换为JSON
	Scan(value any) error                  // 从数据库查询解析
	Value() (driver.Value, error)          // 转换为数据库存储
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
