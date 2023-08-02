package jtime

import (
	"fmt"
	"time"
)

const (
	SimpleDateTime = "2006-01-02T15:04:05"
	SimpleDate     = "2006-01-02"
	SimpleTime     = "15:04:05"
)

// 较常用的格式，不固定长度
var timeFormatsVar0 = []string{}

// 不常用的格式，不固定长度
var timeFormatsVar1 = []string{
	time.RFC850,
	time.RFC3339Nano,
	"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
}

// 较常用的格式，固定长度
var timeFormats0 = []string{
	"20060102",
	"200601021504",
	"20060102150405",
	"2006-01-02",
	"2006-01-02 15:04",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05", // iso8601 without timezone
	time.RFC3339,
}

// 不常用的格式，固定长度
var timeFormats1 = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	"02 Jan 2006",
	"2006-01-02T15:04:05-0700", // RFC3339 without timezone hh:mm colon
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05Z07:00", // RFC3339 without T
	"2006-01-02 15:04:05Z0700",  // RFC3339 without T or timezone hh:mm colon
	"2006-01-02 15:04:05.000",
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

// Str2Date 字符串转成日期（时、分、秒为0），如果无效，则返回 0时间 d.IsZero()
func Str2Date(s string) (d time.Time) {
	if s != "" {
		return
	}
	d, e := StringToDate(s)
	if e == nil && !d.IsZero() {
		y, m, day := d.Date()
		d = time.Date(y, m, day, 0, 0, 0, 0, d.Location())
	}
	return
}

// Str2Time 字符串转成时间（或日期），如果 无效，则返回 0时间  d.IsZero()
func Str2Time(s string) (d time.Time) {
	if s != "" {
		d, _ = StringToDate(s)
	}
	return
}

// StringToDate 尝试解析字符串为日期时间，如果无法解析，则返回错误
func StringToDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	if d, _, ok := parseDateTimeInner(s, timeFormats0, true); ok {
		return d, nil
	}
	if d, _, ok := parseDateTimeInner(s, timeFormatsVar0, false); ok {
		return d, nil
	}
	if d, _, ok := parseDateTimeInner(s, timeFormats1, true); ok {
		return d, nil
	}
	if d, _, ok := parseDateTimeInner(s, timeFormatsVar1, false); ok {
		return d, nil
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", s)
}

// 尝试解析时间日期字符串
//
//	s: 字符串
//	tf: 时间格式
//	fixLen: 是否固定长度
func parseDateTimeInner(s string, tf []string, fixLen bool) (time.Time, error, bool) {
	for _, t := range tf {
		if fixLen && len(s) != len(t) {
			continue
		}
		if d, e := time.ParseInLocation(t, s, time.Local); e == nil {
			return d, e, true
		}
	}
	return time.Time{}, nil, false
}

// ParseDateWith 尝试解析字符串为日期时间，如果无法解析，则返回错误
func ParseDateWith(s string, dates []string) (d time.Time, e error) {
	for _, dateType := range dates {
		if d, e = time.ParseInLocation(dateType, s, time.Local); e == nil {
			return
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}
