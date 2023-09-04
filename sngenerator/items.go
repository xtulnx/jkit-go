package sngenerator

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

const (
	Digit    = "0123456789"
	Alpha    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AlphaNum = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ" // 去掉易混淆的 IO
	Hex      = Digit + "ABCDEF"
)

const (
	AlignLeft  = 0
	AlignRight = 1
)

// TConst 常量
type TConst string

func (t TConst) V(s *Session) (string, error) {
	return string(t), nil
}

// TSession 子层
type TSession struct {
	Exp ExpVal
}

func (t *TSession) V(s *Session) (string, error) {
	return t.Exp.V(s.Clone())
}

// TEnv 环境变量,依赖 FnEnv
type TEnv struct {
	Name ExpVal
}

func (t TEnv) V(s *Session) (string, error) {
	if s.FnEnv == nil {
		return "", nil
	}
	v, e := t.Name.V(s)
	if e != nil {
		return "", e
	}
	return s.FnEnv(s.Ctx, v)
}

// TCode 业务代号
type TCode struct {
	Exp      ExpVal
	OnlyCode bool
	_code    string
}

func (t *TCode) Code(s *Session) (string, error) {
	if t.Exp == nil || t.Exp == t {
		return "", nil
	}
	if c, ok := t.Exp.(ExpCode); ok {
		return c.Code(s)
	}
	return t.Exp.V(s)
}

func (t *TCode) V(s *Session) (string, error) {
	if t.Exp == nil || t.Exp == t || t.OnlyCode {
		return "", nil
	}
	return t.Exp.V(s)
}

// TIncr 计数
type TIncr struct {
	Min    int64  // 最小值
	Step   int64  // 步长, 默认为1
	Format string // 格式化字符串，默认十进制整数
	Len    int    // 0:不限制长度
	Align  int    // 0:左对齐 1:右对齐(默认)
	Pad    string // 补0字符，默认是 "0", 如果长度大于1，则随机填充
}

func (t *TIncr) SetFormat(format string) *TIncr {
	t.Format = format
	return t
}

func (t *TIncr) SetAlign(len int, align int, pad string) *TIncr {
	t.Len = len
	t.Align = align
	t.Pad = pad
	return t
}

// KitFill 填充文本串
// s: 原始串
// l: 填充后的长度
// alignLeft: 是否左对齐,左对齐时填充在右边
// pad: 填充字符
// rnd: 随机数生成器
func KitFill(s string, l int, alignLeft bool, pad string, rnd *rand.Rand) string {
	if l <= len(s) {
		return s
	}
	s2 := make([]byte, l)
	var l1, l2 int
	if alignLeft {
		l1, l2 = 0, len(s)
	} else {
		l1, l2 = l-len(s), 0
	}
	copy(s2[l1:], s)
	if pad == "" {
		pad = "0"
	}
	if len(pad) > 1 {
		for i := 0; i < l-len(s); i++ {
			s2[i+l2] = pad[rnd.Intn(len(pad))]
		}
	} else {
		for i := 0; i < l-len(s); i++ {
			s2[i+l2] = pad[0]
		}
	}
	return string(s2)
}

func (t *TIncr) V(s *Session) (string, error) {
	if s.FnCnt == nil {
		return "", nil
	}
	k := strings.Join(s.Code, ".")
	v, e := s.FnCnt(s.Ctx, k, t.Min, t.Step)
	if e != nil {
		return "", e
	}
	var s1 string
	if t.Format != "" {
		s1 = fmt.Sprintf(t.Format, v)
	} else {
		s1 = strconv.FormatInt(v, 10)
	}
	s1 = KitFill(s1, t.Len, t.Align == AlignLeft, t.Pad, s.Rnd)
	return s1, nil
}

// TIncrLuck 计数，不显示特殊号码 4, 13, 14, 24, 44, 74, 84, 94
type TIncrLuck struct {
	Min   int64  // 最小值
	Len   int    // 0:不限制长度
	Align int    // 0:左对齐 1:右对齐(默认)
	Pad   string // 补0字符，默认是 "0", 如果长度大于1，则随机填充
}

var badNumber = []string{
	"4", "13", "14", "24", "44", "74", "84", "94",
}

func (t *TIncrLuck) V(s *Session) (string, error) {
	// TODO:
	if s.FnCnt == nil {
		return "", nil
	}
	k := strings.Join(s.Code, ".")
	// TODO: 最多尝试10次
	v, e := s.FnCnt(s.Ctx, k, t.Min, 1)
	if e != nil {
		return "", e
	}
	var s1 string
	s1 = strconv.FormatInt(v, 10)

	s1 = KitFill(s1, t.Len, t.Align == AlignLeft, t.Pad, s.Rnd)
	return s1, nil
}

const (
	DateFmtYear  = "2006"
	DateFmtMonth = "200601"
	DateFmtDay   = "20060102"
	DateFmtWeek  = "200601Mon"
)

type TTime struct {
	Format     string // 格式化字符串, 为空时忽略
	IsCode     bool   // 是否是周期
	FormatCode string // 格式化字符串, 为空时使用 Format
}

func (t TTime) Code(s *Session) (string, error) {
	if t.IsCode {
		if t.FormatCode != "" {
			return s.Now.Format(t.FormatCode), nil
		}
		if t.Format != "" {
			return s.Now.Format(t.Format), nil
		}
		return s.Now.Format(DateFmtDay), nil
	}
	return "", nil
}

func (t TTime) V(s *Session) (string, error) {
	if t.Format != "" {
		return s.Now.Format(t.Format), nil
	}
	return "", nil
}

// TRand 随机数
type TRand struct {
	Len int    // 长度
	Chr string // 字符集
}

func (t TRand) V(s *Session) (string, error) {
	if t.Len <= 0 || t.Chr == "" {
		return "", nil
	}
	s1 := make([]byte, t.Len)
	for i := 0; i < t.Len; i++ {
		s1[i] = t.Chr[s.Rnd.Intn(len(t.Chr))]
	}
	return string(s1), nil
}

type TRandFill struct {
	Len   int    // 长度, 长度不足时补充随机数
	Chr   string // 字符集
	Align int    // 0:左对齐(默认) 1:右对齐
	Exp   ExpVal // 填充表达式
}

func (t TRandFill) V(s *Session) (string, error) {
	if t.Len <= 0 || t.Chr == "" || t.Exp == nil {
		return "", nil
	}
	v, e := t.Exp.V(s)
	if e != nil {
		return "", e
	}
	s1 := KitFill(v, t.Len, t.Align != AlignRight, t.Chr, s.Rnd)
	return s1, nil
}

type TJoin struct {
	Sep  string   // 连接字符
	Args []ExpVal // 参数
}

func (t *TJoin) Add(args ...ExpVal) *TJoin {
	t.Args = append(t.Args, args...)
	return t
}

func (t *TJoin) V(s *Session) (string, error) {
	s1 := make([]string, 0, len(t.Args))
	for _, a := range t.Args {
		if v, ok := a.(ExpCode); ok {
			if c, e := v.Code(s); e != nil {
				return "", e
			} else if c != "" {
				s.Code = append(s.Code, c)
			}
		}
	}
	for _, arg := range t.Args {
		if v, e := arg.V(s); e != nil {
			return "", e
		} else if v != "" {
			s1 = append(s1, v)
		}
	}
	return strings.Join(s1, t.Sep), nil
}
