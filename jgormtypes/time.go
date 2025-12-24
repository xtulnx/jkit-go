package jgormtypes

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type JTime time.Duration

func NewTimeFromString(s string) (JTime, error) {
	var t JTime
	err := t.SetFromString(s)
	return t, err
}

func NewTime(hour, min, sec, nsec int) JTime {
	return JTime(
		time.Duration(hour)*time.Hour +
			time.Duration(min)*time.Minute +
			time.Duration(sec)*time.Second +
			time.Duration(nsec)*time.Nanosecond,
	)
}

func (t *JTime) SetFromTime(src time.Time) {
	*t = NewTime(src.Hour(), src.Minute(), src.Second(), src.Nanosecond())
}

func (t *JTime) SetFromString(str string) (err error) {
	var h, m, s, n int
	s1 := strings.Split(str, ":")
	if len(s1) > 0 {
		h, err = strconv.Atoi(s1[0])
	}
	if err == nil && len(s1) > 1 {
		m, err = strconv.Atoi(s1[1])
	}
	if err == nil && len(s1) > 2 {
		s2 := strings.Split(s1[2], ".")
		if err == nil && len(s2) > 0 {
			s, err = strconv.Atoi(s2[0])
		}
		if err == nil && len(s2) > 1 {
			n, err = strconv.Atoi(s2[1])
		}
	}
	*t = NewTime(h, m, s, n)
	return err
}

func (t JTime) String() string {
	if t.IsValid() {
		if nsec := t.nanoseconds(); nsec > 0 {
			return fmt.Sprintf("%02d:%02d:%02d.%09d", t.hours(), t.minutes(), t.seconds(), nsec)
		} else {
			return fmt.Sprintf("%02d:%02d:%02d", t.hours(), t.minutes(), t.seconds())
		}
	}
	return ""
}

func (t JTime) IsValid() bool {
	return time.Duration(t) >= time.Millisecond
}

func (t JTime) IsZero() bool {
	return !t.IsValid()
}

func (t JTime) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *JTime) UnmarshalText(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	err := t.SetFromString(string(b))
	return err
}

func (t *JTime) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return nil
	}
	err = t.SetFromString(strings.Trim(string(b), `"`))
	return err
}

func (t JTime) MarshalJSON() ([]byte, error) {
	s1 := t.String()
	if s1 == "" {
		return []byte("null"), nil
	}
	return bytes.Join([][]byte{{'"'}, []byte(s1), {'"'}}, nil), nil
}

func (t *JTime) Scan(value any) (err error) {
	switch v := value.(type) {
	case []byte:
		err = t.SetFromString(string(v))
	case string:
		err = t.SetFromString(v)
	case time.Time:
		t.SetFromTime(v)
	default:
		err = fmt.Errorf("failed to scan value: %v", v)
	}
	return err
}

func (t JTime) Value() (driver.Value, error) {
	s1 := t.String()
	if s1 == "" {
		return nil, nil
	}
	return s1, nil
}

func (t JTime) hours() int {
	return int(time.Duration(t).Truncate(time.Hour).Hours())
}

func (t JTime) minutes() int {
	return int((time.Duration(t) % time.Hour).Truncate(time.Minute).Minutes())
}

func (t JTime) seconds() int {
	return int((time.Duration(t) % time.Minute).Truncate(time.Second).Seconds())
}

func (t JTime) nanoseconds() int {
	return int((time.Duration(t) % time.Second).Nanoseconds())
}

//////////////////////////////////////////////////////////////////

// JDailyTimeRange 每天内的时间范围
type JDailyTimeRange struct {
	Begin JTime `json:"begin"` // 开始时间点
	End   JTime `json:"end"`   // 结束时间点
}
