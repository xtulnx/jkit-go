package jsonkit

import (
	"github.com/json-iterator/go/extra"
)

func init() {
	// 启用 json 兼容处理 ，
	// 还需要在编译参数加上:  -tags jsoniter
	//  如 go run -tags "jsoniter" main.go
	extra.RegisterFuzzyDecoders()
}
