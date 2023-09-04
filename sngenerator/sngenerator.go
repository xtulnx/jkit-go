package sngenerator

import (
	"context"
	"math/rand"
	"time"
)

// 辅助工具：序列号生成
//
// jason.liao 2023.08.24

/*

基本元素:

1. 常量
	* 字符串拼接
2. 递增序列号
	* 动态的
	* 批量获取提高效率
	* 最小值
	* 格式化的对齐: 长度、对齐方式、补0字符
3. 时间戳
	* 按年、按月、按日、按周
	* 直接格式化使用
4. 随机数
	* 长度
	* 字符集
5. 业务代号
	* 在计数器中作为前缀

运算符:

1. "+"

函数:

1. 日期格式化
2. 数值格式化
3. 随机数生成
4. 字符串截取

*/

// SnGenerator 序列号生成
type SnGenerator interface {
	Next(ctx context.Context, opt ...OptionSession) (string, error)
}

// ExpVal 计算表达式
type ExpVal interface {
	V(s *Session) (string, error)
}

// ExpCode 计数器的 key
type ExpCode interface {
	Code(s *Session) (string, error)
}

// FnCounter 计数器
type FnCounter func(ctx context.Context, code string, min, step int64) (int64, error)

// FnEnv 环境变量
type FnEnv func(ctx context.Context, name string) (string, error)

type snGenerator struct {
	rule  string
	rnd   *rand.Rand
	items ExpVal
	fnEvn FnEnv
	fnCnt FnCounter
}

func NewGenerator(rule string, expr ExpVal, fnEnv FnEnv, fnCnt FnCounter) (SnGenerator, error) {
	s := &snGenerator{
		rule:  rule,
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
		items: expr,
		fnEvn: fnEnv,
		fnCnt: fnCnt,
	}
	return s, nil
}

func (s *snGenerator) Next(ctx context.Context, opts ...OptionSession) (string, error) {
	ss := &Session{
		Ctx:   ctx,
		Now:   time.Now(),
		Code:  nil,
		Rnd:   s.rnd,
		FnEnv: s.fnEvn,
		FnCnt: s.fnCnt,
	}
	for _, o := range opts {
		o(ss)
	}
	return s.items.V(ss)
}

// Session 会话
type Session struct {
	Ctx    context.Context
	Now    time.Time
	Prefix string
	Code   []string
	Rnd    *rand.Rand
	FnEnv  FnEnv
	FnCnt  FnCounter
}

func (s *Session) Clone() *Session {
	return &Session{
		Ctx:    s.Ctx,
		Now:    s.Now,
		Prefix: s.Prefix,
		Code:   nil,
		Rnd:    s.Rnd,
		FnEnv:  s.FnEnv,
		FnCnt:  s.FnCnt,
	}
}

type OptionSession func(s *Session)

func NewMapEnv(m map[string]string) FnEnv {
	return func(ctx context.Context, name string) (string, error) {
		return m[name], nil
	}
}

func NewMapCounter(m map[string]int64) FnCounter {
	return func(ctx context.Context, code string, min, step int64) (int64, error) {
		v, _ := m[code]
		if v < min {
			v = min
		} else {
			v += step
		}
		m[code] = v
		return v, nil
	}
}
