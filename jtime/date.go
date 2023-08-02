package jtime

import "time"

// Age 简单计算年龄（只看年份）
func Age(tBirth time.Time) int {
	if tBirth.IsZero() {
		return 0
	}
	return Age0(tBirth, time.Now())
}

// Age0 简单计算年龄（只看年份）
//
//	tBirth: 出生日期
//	tNow: 当前日期
func Age0(tBirth, tNow time.Time) int {
	if tBirth.IsZero() {
		return 0
	}
	var age = tNow.Year() - tBirth.Year()
	if age < 0 {
		age = 0
	}
	return age
}

// DateFillAB 补全两个日期节点
//
//	a,b 起始日期、结束日期
//	daySize 如果 a、b 一个有效一个为空日期，则以 daySize 计算另外一个日期，
//	day0 如果 a、b 都是空日期，则以「今天-day0」作为开始日期
//
// 例如，
//
//	如果扩展到本周，可以是 DateFillAB(a,b,7,int(time.Now().Weekday())-1) // 周一到周天
//	如果扩展到本月，可以是 DateFillAB(a,b,31,int(time.Now().Day())-1) // 周一到周天
func DateFillAB(a, b time.Time, daySize, day0 int) (time.Time, time.Time) {
	if !a.IsZero() {
		if b.IsZero() {
			return a, a.AddDate(0, 0, daySize)
		} else {
			return a, b
		}
	} else if !b.IsZero() {
		return b.AddDate(0, 0, -daySize), b
	} else {
		a = TruncDay(time.Now()).AddDate(0, 0, -day0)
		return a, a.AddDate(0, 0, daySize)
	}
}

// DateFillMonth 扩展到月内时间
//  1. 如果 a、b 都是空日期，则以「今天」作为开始日期
//  2. 如果 a、b 一个有效一个为空日期，则以另外一个日期的年月作为开始日期
func DateFillMonth(a, b time.Time, minDay int) (time.Time, time.Time) {
	zeroA, zeroB := a.IsZero(), b.IsZero()
	if !zeroA && !zeroB {
		return a, b
	}
	var y int
	var m time.Month
	if !zeroA {
		y, m, _ = a.Date()
	} else if !zeroB {
		y, m, _ = b.Date()
	} else {
		y, m, _ = time.Now().Date()
	}
	if zeroA {
		a = time.Date(y, m, 1, 0, 0, 0, 0, b.Location())
	}
	if zeroB {
		b = time.Date(y, m+1, 1, 0, 0, 0, 0, b.Location())
	}
	if b.Sub(a) < time.Duration(minDay)*time.Hour*24 {
		if !zeroB {
			a = b.AddDate(0, 0, -minDay)
		} else {
			b = a.AddDate(0, 0, minDay)
		}
	}
	return a, b
}
