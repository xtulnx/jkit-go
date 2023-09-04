package jnet

import (
	"github.com/google/shlex"
	"net/http"
	"strings"
)

type CurlCmd struct {
	Url     string
	Headers http.Header
	Data    string
	Method  string
}

func ParseCurl(s string) (*CurlCmd, error) {
	args, err := shlex.Split(strings.ReplaceAll(s, "\\\n", " ")) // 过滤掉 行尾 干扰
	if err != nil {
		return nil, err
	}
	return ParseCurlCmd(args), nil
}

// ParseCurlCmd 解析「简单的」curl
// args,err:= shlex.Split(strings.ReplaceAll(s,"\\\n"," "))
func ParseCurlCmd(args []string) *CurlCmd {
	c := &CurlCmd{Headers: make(http.Header)}
	i, n := 0, len(args)
	fnNext := func() string {
		for ; i < n; i++ {
			s1 := strings.TrimSpace(args[i])
			//if s1 != "" {
			i++
			return s1
			//}
			//return s1
		}
		return ""
	}
	for {
		f1 := fnNext()
		if f1 == "" {
			break
		}
		switch f1 {
		case "curl":
			continue
		case "-X", "--request":
			c.Method = fnNext()
		case "-H", "--header":
			ss := strings.SplitN(fnNext(), ":", 2)
			if len(ss) == 2 {
				c.Headers.Add(strings.TrimSpace(ss[0]), strings.TrimSpace(ss[1]))
			}
		case "--compressed":
			continue
		case "-k", "--insecure":
			continue
		case "-d", "--data", "--data-ascii":
			c.Data = fnNext()
		case "--data-raw":
			c.Data = fnNext()
		case "--url":
			c.Url = fnNext()
		default:
			if f1[0] == '-' {
				_ = fnNext()
			} else {
				c.Url = f1
			}
		}
	}
	return c
}
