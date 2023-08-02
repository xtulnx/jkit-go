package jtime

import "time"

type TruncTime time.Time

func (t0 TruncTime) V() time.Time {
	return time.Time(t0)
}

// Truncate 截断对齐；不支持按月、按年
// 如果是按其他维度，比如 一天内小时、小时内分钟 等，则要先对齐到目标点，再使用差值计算。
func (t0 TruncTime) Truncate(d time.Duration) TruncTime {
	return TruncTime(Truncate(time.Time(t0), d))
}

// Truncate2 依次对齐
func (t0 TruncTime) Truncate2(d1, d2 time.Duration) TruncTime {
	return TruncTime(Truncate2(time.Time(t0), d1, d2))
}

// TruncDay 按天对齐
func (t0 TruncTime) TruncDay() TruncTime {
	return TruncTime(TruncDay(time.Time(t0)))
}

// NextDayStart 下一天的开始
func (t0 TruncTime) NextDayStart() TruncTime {
	return TruncTime(NextDayStart(time.Time(t0)))
}

// TruncMonth 当月开始
func (t0 TruncTime) TruncMonth() TruncTime {
	return TruncTime(TruncMonth(time.Time(t0)))
}

// NextMonthStart 次月开始
func (t0 TruncTime) NextMonthStart() TruncTime {
	return TruncTime(NextMonthStart(time.Time(t0)))
}

func (t0 TruncTime) Add(y, m, d int) TruncTime {
	return TruncTime(time.Time(t0).AddDate(y, m, d))
}

func (t0 TruncTime) TruncWeek() TruncTime {
	return TruncTime(TruncWeek(time.Time(t0)))
}

func (t0 TruncTime) TruncQuarter() TruncTime {
	return TruncTime(TruncQuarter(time.Time(t0)))
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// TruncDay 按天截断
func TruncDay(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// Truncate 截断对齐；不支持按月、按年
//
//	如果是按其他维度，比如 一天内小时、小时内分钟 等，则要先对齐到目标点，再使用差值计算。
//
// 用法:
//
//	t1 => 2021-06-17 18:32:11 +0800 CST
//	按天: TimeTruncate(t1, time.Hour*24) => 2021-06-17 00:00:00 +0800 CST
//	按7小时: TimeTruncate(t1, time.Hour*7) => 2021-06-17 12:00:00 +0800 CST
//	按7小时（当天）: TimeTruncate2(t1, time.Hour*24, time.Hour*7) => 2021-06-17 14:00:00 +0800 CST
//	按小时: TimeTruncate(t1, time.Hour) => 2021-06-17 18:00:00 +0800 CST
//	按分钟: TimeTruncate(t1, time.Minute) => 2021-06-17 18:32:00 +0800 CST
func Truncate(t0 time.Time, d time.Duration) time.Time {
	if t0.IsZero() {
		return t0
	}
	_, dif := t0.Zone()
	addDiff := time.Second * time.Duration(dif)
	return t0.Add(addDiff).Truncate(d).Add(-addDiff)
}

// Truncate2 依次对齐
func Truncate2(t0 time.Time, d1, d2 time.Duration) time.Time {
	if t0.IsZero() {
		return t0
	}
	_, dif := t0.Zone()
	addDiff := time.Second * time.Duration(dif)
	t1 := t0.Add(addDiff).Truncate(d1).Add(-addDiff)
	return t1.Add(t0.Sub(t1).Truncate(d2))
}

// NextDayStart 下一天的开始
func NextDayStart(t0 time.Time) time.Time {
	if t0.IsZero() {
		return t0
	}
	y, m, d := t0.Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, t0.Location())
}

// TruncMonth 取月初日期
func TruncMonth(t0 time.Time) time.Time {
	if t0.IsZero() {
		return t0
	}
	y, m, _ := t0.Date()
	return time.Date(y, m, 0, 0, 0, 0, 0, t0.Location())
}

// NextMonthStart 下个月的开始
func NextMonthStart(t0 time.Time) time.Time {
	if t0.IsZero() {
		return t0
	}
	y, m, _ := t0.Date()
	return time.Date(y, m+1, 1, 0, 0, 0, 0, t0.Location())
}

// MonthAdd 月份加减，如果是无效时间则跳过
func MonthAdd(t0 time.Time, inc int) time.Time {
	if t0.IsZero() {
		return t0
	}
	return t0.AddDate(0, inc, 0)
}

// TruncWeek 按周对齐
func TruncWeek(t0 time.Time) time.Time {
	if t0.IsZero() {
		return t0
	}
	y, m, d := t0.Date()
	w := t0.Weekday()
	return time.Date(y, m, d-int(w), 0, 0, 0, 0, t0.Location())
}

// TruncQuarter 按季度对齐
func TruncQuarter(t0 time.Time) time.Time {
	if t0.IsZero() {
		return t0
	}
	y, m, _ := t0.Date()
	quarter := (m-1)/3 + 1
	return time.Date(y, (quarter-1)*3+1, 1, 0, 0, 0, 0, t0.Location())
}
