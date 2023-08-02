package jtime

import "time"

var timeFormatsForExpandSection = [][]string{
	// 按年
	{"2006", "2006年"},
	// 按月
	{"200601", "2006-01", "Jan 2006", "2006年01月"},
	// 按天
	{"20060102", "2006-01-02", "02 Jan 2006", "2006年01月02日"},
	// 时
	{"2006010215", "2006-01-02 15", "02 Jan 2006", "2006年01月02日"},
	// 分
	{"200601021504", "2006-01-02 15:04", "2006-01-02T15:04"},
	// 秒，其他
	{
		"20060102150405",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05", // iso8601 without timezone
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		time.ANSIC,
		time.UnixDate,
		"2006-01-02 15:04:05.999999999 -0700 MST", // Time.String()
		"2006-01-02T15:04:05-0700",                // RFC3339 without timezone hh:mm colon
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
	},
}

// Str2TimeExpand 扩展包含当前的时段；如参数是按天，则扩展到明天；如果参数精确到小时，则扩展到下一小时
func Str2TimeExpand(s string) (d time.Time) {
	for i := range timeFormatsForExpandSection {
		for _, dateType := range timeFormatsForExpandSection[i] {
			if t1, e1 := time.ParseInLocation(dateType, s, time.Local); e1 == nil {
				year, month, day := t1.Date()
				hour, min, sec := t1.Clock()
				switch i {
				case 0:
					d = time.Date(year+1, 0, 0, 0, 0, 0, 0, t1.Location())
				case 1:
					d = time.Date(year, month+1, 0, 0, 0, 0, 0, t1.Location())
				case 2:
					d = time.Date(year, month, day+1, 0, 0, 0, 0, t1.Location())
				case 3:
					d = time.Date(year, month, day, hour+1, 0, 0, 0, t1.Location())
				case 4:
					d = time.Date(year, month, day, hour, min+1, 0, 0, t1.Location())
				case 5:
					d = time.Date(year, month, day, hour, min, sec+1, 0, t1.Location())
				default:
					d = t1
				}
				return d
			}
		}
	}
	return Str2Time(s)
}
