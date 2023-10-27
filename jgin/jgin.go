package jgin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
)

// gin 框架的辅助函数

var ginBinding = map[string]binding.Binding{}
var bindingHandler []BindingHandler

// RegBinding 注册请求与 model 的解析器
func RegBinding(mime string, b binding.Binding) {
	ginBinding[mime] = b
}

type BindingHandler func(context2 *gin.Context, req any) error

func RegBindingHandler(h BindingHandler) {
	bindingHandler = append(bindingHandler, h)
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

type tWithBinder interface {
	BindGinContext(c *gin.Context) error
}

type tWithNow interface {
	SetNow(now time.Time)
}
type tWithIP interface {
	SetIP(string)
}
type tWithCtx interface {
	SetCtx(ctx context.Context)
	SetCtxValue(k, v interface{})
}

// GinMustBind 示例：解析请求参数
func GinMustBind(c *gin.Context, obj interface{}) error {
	reqMethod, reqContentType := c.Request.Method, c.ContentType()
	var b binding.Binding = nil
	if ginBinding != nil {
		if reqMethod == http.MethodGet {
			b, _ = ginBinding[binding.MIMEPOSTForm]
		} else {
			b, _ = ginBinding[reqContentType]
		}
	}
	if b == nil {
		b = binding.Default(reqMethod, reqContentType)
	}
	err := c.MustBindWith(obj, b)
	if err != nil {
		return err
	}
	if m, ok := obj.(tWithNow); ok {
		m.SetNow(time.Now())
	}
	if m, ok := obj.(tWithCtx); ok {
		m.SetCtx(c.Request.Context())
	}
	if m, ok := obj.(tWithIP); ok {
		m.SetIP(c.ClientIP())
	}
	if m, ok := obj.(tWithBinder); ok {
		err = m.BindGinContext(c)
		if err != nil {
			return err
		}
	}
	if bindingHandler != nil {
		for _, h := range bindingHandler {
			err = h(c, obj)
			if err != nil {
				return err
			}
		}
	}
	return err
}
