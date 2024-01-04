package jgormtypes

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type JTime time.Duration

func (ct *JTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return nil
	} else {
		var hour, min, sec, nsec int
		_, _ = fmt.Sscanf(s, "%02d:%02d:%02d.%09d", &hour, &min, &sec, &nsec)
		t := time.Duration(hour)*time.Hour +
			time.Duration(min)*time.Minute +
			time.Duration(sec)*time.Second +
			time.Duration(nsec)*time.Nanosecond
		*ct = JTime(t)
	}
	return
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

func (ct JTime) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", ct.hours(), ct.minutes(), ct.seconds())
}

func (ct JTime) MarshalJSON() ([]byte, error) {
	if ct == 0 {
		return []byte("null"), nil
	}
	return bytes.Join([][]byte{{'"'}, []byte(ct.String()), {'"'}}, nil), nil
}

//////////////////////////////////////////////////////////////////

// JDailyTimeRange 每天内的时间范围
type JDailyTimeRange struct {
	Begin JTime `json:"begin"` // 开始时间点
	End   JTime `json:"end"`   // 结束时间点
}
