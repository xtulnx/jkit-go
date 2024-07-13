package jtypes

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
)

func TestTypes(t *testing.T) {
	typesPrice(t)
}

func typesPrice(t *testing.T) {
	data := []float64{
		1234567890123456789.12345678,
		1234567890123456789.123,
		1234567890123456789.00001234,
		1234567890123456789.000,
		1234567890123456.12345678,
		1234567890123456.123,
		1234567890123456.00005678,
		1234567890123456.00000678,
		12345678901234.12345678,
		12345678901234.00005678,
		12345678901234.00000678,
		12345678901234.00000078,
		12345678901234.123,
		12345678901234.000,
		1234567890.12345678,
		1234567890.00000678,
		1234567890.00005678,
		1234567890.1200,
		1234567890.00000,
		120.12345678,
		0.12345678,
		0.00345678,
		0.00005678,
		0.00000678,
	}
	for _, i := range data {
		b1, _ := json.Marshal(i)
		t.Log(i, "=> FMT:", JPrice(i).String(), ", JSON:", string(b1))
	}
}

func TestPrice2(t *testing.T) {
	fnS := func(s string) JPrice {
		p, _ := NewPriceFromString(s)
		return p
	}
	fnStr := func(v JPrice) string {
		return fmt.Sprintf("float=%v, str=%s, cent=%v, fixed=%v, round=%v", v.Float(), v.String(), v.ToCent(), v.ToFixed(), v.Round().Float())
	}
	fnEq := func(inc JPrice, ok bool) func(v JPrice) string {
		return func(v JPrice) string {
			a := inc + v
			e := a.Equal(v)
			if e != ok {
				t.Errorf("v=%v, a=%v, a.Equal(v)=%v, want %v", v, a, e, ok)
			}
			return fmt.Sprintf("%s.Equal(%s) => %v", v.String(), a.String(), e)
		}
	}
	for _, v := range []struct {
		N    string
		P    JPrice
		F    func(v JPrice) string
		Want string
	}{
		{"1/3", 1.0 / 3, fnStr, "float=0.3333333333333333, str=0.33, cent=33, fixed=3333, round=0.3333"},
		{"2/3", 2.0 / 3, fnStr, "float=0.6666666666666666, str=0.67, cent=67, fixed=6667, round=0.6667"},
		{"2024/31", 2024.0 / 31, fnStr, "float=65.29032258064517, str=65.29, cent=6529, fixed=652903, round=65.2903"},
		{"-1/3", -1.0 / 3, fnStr, "float=-0.3333333333333333, str=-0.33, cent=-33, fixed=-3333, round=-0.3333"},
		{"-2/3", -2.0 / 3, fnStr, "float=-0.6666666666666666, str=-0.67, cent=-67, fixed=-6667, round=-0.6667"},
		{"2/3-0.002", 2.0/3 - 0.002, fnStr, "float=0.6646666666666666, str=0.66, cent=66, fixed=6647, round=0.6647"},
		{"2/3-0.0002", 2.0/3 - 0.0002, fnStr, "float=0.6664666666666667, str=0.67, cent=67, fixed=6665, round=0.6665"},
		{"Nan", fnS("Nan"), fnStr, "float=NaN, str=NaN, cent=-9223372036854775808, fixed=-9223372036854775808, round=NaN"},
		{"678,901 234,567.890123456", fnS("678,901 234,567.890123456"), fnStr, "float=6.789012345678901e+11, str=678901234567.89, cent=67890123456789, fixed=6789012345678901, round=6.789012345678901e+11"},
		{"2/3 - 0.002 notEqual", 2.0 / 3, fnEq(-0.02, false), "0.67.Equal(0.65) => false"},
		{"2/3 - 0.002 notEqual", 2.0 / 3, fnEq(-0.002, false), "0.67.Equal(0.66) => false"},
		{"2/3 - 0.0002 equal", 2.0 / 3, fnEq(-0.0002, true), "0.67.Equal(0.67) => true"},
		{"2/3 - 0.00002 equal", 2.0 / 3, fnEq(-0.00002, true), "0.67.Equal(0.67) => true"},
		{"1234567890 - 7890 notEqual", 1234567890, fnEq(-7890, false), "1234567890.Equal(1234560000) => false"},
		{"pi*1e20", math.Pi * 1e8, fnStr, "float=3.1415926535897934e+08, str=314159265.36, cent=31415926536, fixed=3141592653590, round=3.14159265359e+08"},
		// 9,223,372,036,854,775,807
		{"max cent", NewPriceFromCent(1<<63 - 1), fnStr, "float=9.223372036854776e+16, str=92233720368547760, cent=-9223372036854775808, fixed=-9223372036854775808, round=9.223372036854776e+16"},
		{"max fixed", NewPriceFromFixed(1<<63 - 1), fnStr, "float=9.223372036854776e+14, str=922337203685477.62, cent=92233720368547760, fixed=-9223372036854775808, round=9.223372036854776e+14"},
		// 4,611,686,018,427,387,903
		{"cent62", NewPriceFromCent(1<<62 - 1), fnStr, "float=4.611686018427388e+16, str=46116860184273880, cent=4611686018427387904, fixed=-9223372036854775808, round=4.611686018427388e+16"},
		{"fixed62", NewPriceFromFixed(1<<62 - 1), fnStr, "float=4.611686018427388e+14, str=461168601842738.81, cent=46116860184273880, fixed=4611686018427387904, round=4.611686018427388e+14"},
	} {
		v1 := v
		t.Run(v.N, func(t *testing.T) {
			s1 := v1.F(v1.P)
			if s1 != v1.Want {
				t.Errorf("%s => %#v, [%s], want [%s]", v1.N, v1.P, s1, v1.Want)
			}
		})
	}
}
