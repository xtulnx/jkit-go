package jtime

import (
	"testing"
	"time"
)

func TestAge(t *testing.T) {
	tests := []struct {
		args string
		want int
	}{
		{"", 0},
		{"1990-01-01", 33},
		{"1990-06-01", 33},
		{"1990-09-01", 33},
		{"1990-08-01", 33},
		{"1990-08-02", 33},
		{"1990-12-01", 33},
	}
	tNow, _ := time.Parse("2006-01-02", "2023-08-02")
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			tBirth, _ := time.Parse("2006-01-02", tt.args)
			if got := Age0(tBirth, tNow); got != tt.want {
				t.Errorf("Age() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateFill(t *testing.T) {
	tests := []struct {
		name          string
		a, b          string
		daySize, day0 int
		a1, b1        string
	}{
		{"7天,对齐到周一开始", "", "", 7, int(time.Now().Weekday() - 1), "", ""},
		{"31天,对齐到周一开始", "", "", 31, int(time.Now().Weekday() - 1), "", ""},
		{"缺失右端日期，按天数补上", "2023-07-06", "", 3, 1, "2023-07-06", "2023-07-09"},             //
		{"day0只用在 a,b 都是空的情况，否则不使用", "2023-07-06", "", 3, 3, "2023-07-06", "2023-07-09"}, //
		{"", "2023-07-06", "", 30, 3, "2023-07-06", "2023-08-05"},
		{"缺失左端日期，按天数补上", "", "2023-07-06", 30, 3, "2023-06-06", "2023-07-06"},
		{"都有值时忽略 daySize,day0", "2023-06-20", "2023-07-06", 30, 3, "2023-06-20", "2023-07-06"},
	}
	for _, tt := range tests {
		name := tt.name
		if name == "" {
			name = tt.a + " " + tt.b
		}
		t.Run(name, func(t *testing.T) {
			a, _ := time.Parse("2006-01-02", tt.a)
			b, _ := time.Parse("2006-01-02", tt.b)
			gotA, gotB := DateFillAB(a, b, tt.daySize, tt.day0)
			a1, b1 := gotA.Format("2006-01-02"), gotB.Format("2006-01-02")
			if tt.a1 != "" && tt.b1 != "" {
				if tt.a1 != a1 || tt.b1 != b1 {
					t.Errorf("DateFillAB() = %s,%s, want %s,%s", a1, b1, tt.a1, tt.b1)
				}
			} else {
				t.Logf("gotA=%v[%v], gotB=%v[%v]", a1, gotA.Weekday(), b1, gotB.Weekday())
			}
		})
	}
}

func TestDateFillMonth(t *testing.T) {
	test := []struct {
		name   string
		a, b   string
		minDay int
		a1, b1 string
	}{
		{"a,b 都有则直接返回", "2023-07-06", "2023-07-10", 30, "2023-07-06", "2023-07-10"},
		{"缺失a，天数足够，则对齐到月初", "", "2023-07-20", 12, "2023-07-01", "2023-07-20"},
		{"缺失a，天数不足够，则按天数计算", "", "2023-07-20", 22, "2023-06-28", "2023-07-20"},
		{"缺失b，天数足够，则对齐到下个月初", "2023-07-20", "", 8, "2023-07-20", "2023-08-01"},
		{"缺失b，天数不足够，则按天数计算", "2023-07-20", "", 22, "2023-07-20", "2023-08-11"},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := time.Parse("2006-01-02", tt.a)
			b, _ := time.Parse("2006-01-02", tt.b)
			gotA, gotB := DateFillMonth(a, b, tt.minDay)
			a1, b1 := gotA.Format("2006-01-02"), gotB.Format("2006-01-02")
			if tt.a1 != "" || tt.b1 != "" {
				if tt.a1 != a1 || tt.b1 != b1 {
					t.Errorf("DateFillMonth() = %s,%s, want %s,%s", a1, b1, tt.a1, tt.b1)
				}
			} else {
				t.Logf("gotA=%v[%v], gotB=%v[%v]", a1, gotA.Weekday(), b1, gotB.Weekday())
			}
		})
	}

}
