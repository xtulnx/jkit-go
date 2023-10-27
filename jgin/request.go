package jgin

import (
	"context"
	"time"
)

// 封装一些请求相关的定义

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// ReqWithNow 统一业务时间点
type ReqWithNow struct {
	now time.Time
}

func (R *ReqWithNow) SetNow(now time.Time) {
	R.now = now
}
func (R *ReqWithNow) GetNow() time.Time {
	if R.now.IsZero() {
		R.now = time.Now()
	}
	return R.now
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// ReqWithIP 客户端 ip
type ReqWithIP struct {
	ip string
}

func (R *ReqWithIP) SetIP(ip string) {
	R.ip = ip
}
func (R *ReqWithIP) GetIP() string {
	return R.ip
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// ReqWithCtx 上下文环境
type ReqWithCtx struct {
	ctx context.Context
}

func (R *ReqWithCtx) SetCtx(ctx context.Context) {
	R.ctx = ctx
}
func (R *ReqWithCtx) GetCtx() context.Context {
	return R.ctx
}
func (R *ReqWithCtx) SetCtxValue(k, v interface{}) {
	R.ctx = context.WithValue(R.ctx, k, v)
}
func (R *ReqWithCtx) GetCtxValue(k interface{}) interface{} {
	return R.ctx.Value(k)
}
