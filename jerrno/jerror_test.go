package jerrno

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {

	var e0 ErrEx = NewErr("Error0")
	var e1 ErrEx = NewErrWithCode(1, "Error1")
	var e2 ErrEx = NewErrWithHttpCode(200, 2, "Error2")
	fmt.Printf("%#v\n", e0)
	fmt.Printf("%#v\n", e1)
	fmt.Printf("%#v\n", e2)
	fmt.Printf("%#v\n", Forbidden)
	fmt.Printf("%#v\n", Forbidden.WithMsg("您无权访问，请重新登录或联系管理员"))
	fmt.Printf("%#v\n", e2.WithMsg("New Error2"))
}
