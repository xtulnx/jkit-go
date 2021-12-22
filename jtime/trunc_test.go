package jtime

import (
	"testing"
	"time"
)

func TestTrunc(t *testing.T) {
	for i, t1 := range []time.Time{
		Str2Time("2021-06-17 18:32:11 +0800 CST"), time.Now(),
	} {
		t.Logf("序号: %d", i)
		t.Logf("时间: %s", t1)
		t2 := TruncTime(t1)
		t.Logf("按天: %s", t2.Truncate(time.Hour*24))
		t.Logf("按7小时: %s", t2.Truncate(time.Hour*7))
		t.Logf("按7小时（当天）: %s", t2.Truncate2(time.Hour*24, time.Hour*7))
		t.Logf("按小时: %s", t2.Truncate(time.Hour))
		t.Logf("按分钟: %s", t2.Truncate(time.Minute))
		t.Log("++++++++++++++++")
		t.Logf("当天开始时间: %s", t2.TruncDay())
		t.Logf("次日开始时间: %s", t2.NextDayStart())
		t.Logf("当月开始时间: %s", t2.TruncMonth())
		t.Logf("次月开始时间: %s", t2.NextMonthStart())
		t.Log("----------------\n")
	}

}
