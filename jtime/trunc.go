package jtime

import "time"

type TruncTime time.Time

// Truncate 截断对齐；不支持按月、按年
// 如果是按其他维度，比如 一天内小时、小时内分钟 等，则要先对齐到目标点，再使用差值计算。
func (t0 TruncTime) Truncate(d time.Duration) time.Time {
	_t0 := time.Time(t0)
	_, dif := _t0.Zone()
	addDiff := time.Second * time.Duration(dif)
	return _t0.Add(addDiff).Truncate(d).Add(-addDiff)
}

// Truncate2 依次对齐
func (t0 TruncTime) Truncate2(d1, d2 time.Duration) time.Time {
	_t0 := time.Time(t0)
	_, dif := _t0.Zone()
	addDiff := time.Second * time.Duration(dif)
	t1 := _t0.Add(addDiff).Truncate(d1).Add(-addDiff)
	return t1.Add(_t0.Sub(t1).Truncate(d2))
}

// TruncDay 按天对齐
func (t0 TruncTime) TruncDay() time.Time {
	_t0 := time.Time(t0)
	y, m, d := _t0.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// NextDayStart 下一天的开始
func (t0 TruncTime) NextDayStart() time.Time {
	_t0 := time.Time(t0)
	y, m, d := _t0.Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, time.Local)
}

// TruncMonth 当月开始
func (t0 TruncTime) TruncMonth() time.Time {
	_t0 := time.Time(t0)
	y, m, _ := _t0.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
}

// NextMonthStart 次月开始
func (t0 TruncTime) NextMonthStart() time.Time {
	_t0 := time.Time(t0)
	y, m, _ := _t0.Date()
	return time.Date(y, m+1, 1, 0, 0, 0, 0, time.Local)
}
