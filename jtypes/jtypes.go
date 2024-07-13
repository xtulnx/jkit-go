package jtypes

import (
	"bytes"
	"math"
	"strconv"
)

// 价格，或财务金额，精度4位小数
type JPrice float64

func NewPrice(p float64) JPrice {
	return JPrice(p)
}

func NewPriceFromCent(cent int64) JPrice {
	return JPrice(float64(cent) / 100)
}

func (p JPrice) ToCent() int64 {
	return int64(math.Round(float64(p) * 100))
}

func NewPriceFromFixed(fix int64) JPrice {
	return JPrice(float64(fix) / 1e4)
}

func (p JPrice) ToFixed() int64 {
	return int64(math.Round(float64(p) * 1e4))
}

func (p JPrice) Round() JPrice {
	return JPrice(math.Round(float64(p)*1e4) / 1e4)
}

func NewPriceFromString(s string) (JPrice, error) {
	var p JPrice
	err := p.UnmarshalText([]byte(s))
	return p, err
}

func (p JPrice) Float() float64 {
	return float64(p)
}

const JPriceSize = 2

func (p JPrice) String() string {
	b1 := strconv.AppendFloat(nil, float64(p), 'f', JPriceSize, 64)
	l, i := len(b1), 1
	for i <= JPriceSize && b1[l-i] == '0' {
		i++
	}
	if i > JPriceSize && b1[l-i] == '.' {
		i++
	}
	if i > 1 {
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

	top := 0
	for i, c := range data {
		switch c {
		case ' ', '\t', ',':
			continue
		default:
			if top != i {
				data[top] = data[i]
			}
			top++
		}
	}
	data = data[:top]

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

func JPriceEqual(p1, p2 float64) bool {
	return JPrice(p1).Equal(JPrice(p2))
}

func JPriceNotEqual(p1, p2 float64) bool {
	return JPrice(p1).NotEqual(JPrice(p2))
}

// Equal 相等
func (p JPrice) Equal(p2 JPrice) bool {
	return p.ToCent() == p2.ToCent()
}

// NotEqual 不相等
func (p JPrice) NotEqual(p2 JPrice) bool {
	return !p.Equal(p2)
}
