package jtypes

import (
	"bytes"
	"strconv"
)

// 价格，或财务金额，精度2位小数
type JPrice float64

func (p JPrice) Float() float64 {
	return float64(p)
}

const JPriceSize = 2

func (p JPrice) String() string {
	b1 := strconv.AppendFloat(nil, float64(p), 'f', JPriceSize, 64)
	if JPriceSize > 0 {
		l, i := len(b1), 1
		for i <= JPriceSize && b1[l-i] == '0' {
			i++
		}
		if i > JPriceSize && b1[l-i] == '.' {
			i++
		}
		b1 = b1[:l-i+1]
	}
	return string(b1)
}

func (p JPrice) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (dp *JPrice) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	data = bytes.Trim(data, "\"")
	//data = bytes.Trim(data, "'")
	if len(data) == 0 {
		return nil
	}
	t1, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}
	*dp = JPrice(t1)
	return err
}

func (dp *JPrice) UnmarshalJSON(b []byte) (err error) {
	return dp.UnmarshalText(b)
}

func (p JPrice) MarshalJSON() ([]byte, error) {
	s1 := p.String()
	return []byte(s1), nil
}
