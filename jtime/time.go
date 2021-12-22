package jtime

import (
	"fmt"
	"time"
)

var timeFormats = []string{
	"20060102",
	"200601021504",
	"20060102150405",
	"2006-01-02 15:04",
	"2006-01-02 15:04:05",
	time.RFC3339,
	"2006-01-02T15:04:05", // iso8601 without timezone
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RFC850,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
	"2006-01-02",
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

// StringToDate attempts to parse a string into a time.Time type using a
// predefined list of formats.  If no suitable format is found, an error is
// returned.
func StringToDate(s string) (time.Time, error) {
	return ParseDateWith(s, timeFormats)
}

func ParseDateWith(s string, dates []string) (d time.Time, e error) {
	for _, dateType := range dates {
		if d, e = time.ParseInLocation(dateType, s, time.Local); e == nil {
			return
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}

// Str2Time 字符串转成时间（或日期），如果 无效，则返回 0时间  d.IsZero()
func Str2Time(s string) (d time.Time) {
	var e error
	if s == "" {
		return d
	}
	for _, dateType := range timeFormats {
		if d, e = time.ParseInLocation(dateType, s, time.Local); e == nil {
			return
		}
	}
	return d
}
