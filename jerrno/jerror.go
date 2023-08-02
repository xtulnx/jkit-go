package jerrno

// 统一业务错误码。简单的。
//
// jason.liao 2023.07.29

type ErrEx interface {
	error
	CombineError(e error) ErrEx                // 合并错误信息，如果是同类型的错误则直接返回
	CombineErrorMsg(msg string, e error) ErrEx // 合并错误信息，如果是同类型的错误，只保留最后一个错误
	WithError(e error) ErrEx                   // 附加错误信息
	WithCode(code int) ErrEx                   // 附加错误码
	WithMsg(msg string) ErrEx                  // 附加错误信息
	WithMsgAndError(msg string, e error) ErrEx // 附加错误信息和错误
	IsSame(e error) bool                       // 判断是否是同类型的错误，只按低级别判断
}

// 这里定义一些常用异常，可以直接使用
// 也可以自行 NewErrWithCode(0, "OK") 定义

var (
	OK                  = NewErrWithCode(0, "OK")     //
	BadRequest          = NewErrWithCode(400, "参数有误") //
	Unauthorized        = NewErrWithCode(401, "请先登录")
	Forbidden           = NewErrWithCode(403, "权限不足")
	Conflict            = NewErrWithCode(409, "记录已经存在")
	TooManyRequests     = NewErrWithCode(429, "访问太快，请稍候再试") //
	InternalServerError = NewErrWithCode(500, "服务器错误")      //
	QueryNotFound       = NewErrWithCode(550, "记录并不存在")
	QueryFailed         = NewErrWithCode(551, "查询记录失败")
	ConvertDataFailed   = NewErrWithCode(552, "数据格式有误")
	FailedUpdate        = NewErrWithCode(553, "数据更新失败")
)

func NewErr(msg string) ErrEx {
	return &error0{Msg: msg}
}

// NewErrWithError 带信息、附加异常
//
//	msg: 信息
//	e: 附加异常
func NewErrWithError(msg string, e error) ErrEx {
	return &error0{Msg: msg, Err: e}
}

// NewErrWithCode 带状态码、信息
//
//	code: 状态码
//	msg: 错误信息
func NewErrWithCode(code int, msg string) ErrEx {
	return &error1{Code: code, error0: error0{Msg: msg}}
}

// NewErrWithCodeAndError 带状态码、信息、附加异常
//
//	code: 状态码
//	msg: 信息
//	e: 附加异常
func NewErrWithCodeAndError(code int, msg string, e error) ErrEx {
	return &error1{Code: code, error0: error0{Msg: msg, Err: e}}
}

// NewErrWithHttpCode 创建带http码、状态码、错误信息的错误
//
//	httpCode: http状态码
//	code: 状态码
//	msg: 错误信息
func NewErrWithHttpCode(httpCode, code int, msg string) ErrEx {
	return &error2{HttpCode: httpCode, error1: error1{Code: code, error0: error0{Msg: msg}}}
}

type error0 struct {
	Msg string
	Err error
}

type error1 struct {
	error0
	Code int
}
type error2 struct {
	error1
	HttpCode int
}

func (e error0) Unwrap() error {
	return e.Err
}
func (e error0) Error() string {
	return e.Msg
}
func (e error1) ErrorCode() int {
	return e.Code
}
func (e error2) ErrorHttpCode() int {
	return e.HttpCode
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

func isInnerError(e1 error) ErrEx {
	switch ex := e1.(type) {
	case *error0:
		return ex
	case *error1:
		return ex
	case *error2:
		return ex
	}
	return nil
}

func (e error0) CombineError(e1 error) ErrEx {
	if ex := isInnerError(e1); ex != nil {
		return ex
	}
	return &error0{Msg: e.Msg, Err: e1}
}

func (e error0) CombineErrorMsg(msg string, e1 error) ErrEx {
	if ex := isInnerError(e1); ex != nil {
		return ex
	}
	return &error0{Msg: msg, Err: e1}
}

func (e error0) WithError(e1 error) ErrEx {
	return &error0{Msg: e.Msg, Err: e1}
}

func (e error0) WithCode(code int) ErrEx {
	return &error1{Code: code, error0: error0{Msg: e.Msg, Err: e.Err}}
}

func (e error0) WithMsg(msg string) ErrEx {
	return &error0{Msg: msg, Err: e.Err}
}

func (e error0) WithMsgAndError(msg string, e1 error) ErrEx {
	return &error0{Msg: msg, Err: e1}
}

func (e error0) IsSame(e1 error) bool {
	if ex, ok := e1.(*error0); ok {
		return ex.Msg == e.Msg
	}
	return false
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

func (e error1) CombineError(e1 error) ErrEx {
	if ex := isInnerError(e1); ex != nil {
		return ex
	}
	return &error1{Code: e.Code, error0: error0{Msg: e.Msg, Err: e1}}
}

func (e error1) CombineErrorMsg(msg string, e1 error) ErrEx {
	if ex := isInnerError(e1); ex != nil {
		return ex
	}
	return &error1{Code: e.Code, error0: error0{Msg: msg, Err: e1}}
}

func (e error1) WithError(e1 error) ErrEx {
	return &error1{Code: e.Code, error0: error0{Msg: e.Msg, Err: e1}}
}

func (e error1) WithCode(code int) ErrEx {
	return &error1{Code: code, error0: error0{Msg: e.Msg, Err: e.Err}}
}

func (e error1) WithMsg(msg string) ErrEx {
	return &error1{Code: e.Code, error0: error0{Msg: msg, Err: e.Err}}
}

func (e error1) WithMsgAndError(msg string, e1 error) ErrEx {
	return &error1{Code: e.Code, error0: error0{Msg: msg, Err: e1}}
}

func (e error1) IsSame(e1 error) bool {
	if ex, ok := e1.(*error1); ok {
		return ex.Code == e.Code
	}
	return false
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

func (e error2) CombineError(e1 error) ErrEx {
	if ex := isInnerError(e1); ex != nil {
		return ex
	}
	return &error2{HttpCode: e.HttpCode, error1: error1{Code: e.Code, error0: error0{Msg: e.Msg, Err: e1}}}
}

func (e error2) CombineErrorMsg(msg string, e1 error) ErrEx {
	if ex := isInnerError(e1); ex != nil {
		return ex
	}
	return &error2{HttpCode: e.HttpCode, error1: error1{Code: e.Code, error0: error0{Msg: msg, Err: e1}}}
}

func (e error2) WithError(e1 error) ErrEx {
	return &error2{HttpCode: e.HttpCode, error1: error1{Code: e.Code, error0: error0{Msg: e.Msg, Err: e1}}}
}

func (e error2) WithCode(code int) ErrEx {
	return &error2{HttpCode: e.HttpCode, error1: error1{Code: code, error0: error0{Msg: e.Msg, Err: e.Err}}}
}

func (e error2) WithMsg(msg string) ErrEx {
	return &error2{HttpCode: e.HttpCode, error1: error1{Code: e.Code, error0: error0{Msg: msg, Err: e.Err}}}
}

func (e error2) WithMsgAndError(msg string, e1 error) ErrEx {
	return &error2{HttpCode: e.HttpCode, error1: error1{Code: e.Code, error0: error0{Msg: msg, Err: e1}}}
}

func (e error2) IsSame(e1 error) bool {
	if ex, ok := e1.(*error2); ok {
		return ex.Code == e.Code && ex.HttpCode == e.HttpCode
	}
	return false
}
