package sngenerator

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

//

//var TEST_RULE = []string{
//	"D+序列号{20}",
//	"ST+序列号(1000)",
//}

func Test序列号(t *testing.T) {
	map_cnt := map[string]int64{}
	exper := B.Join(".",
		B.CodeOnly(B.Const("ST")),
		B.Const("90"),
		B.Code(B.Const("ORDER")),
		B.Incr(100000, 39999).SetAlign(8, 1, "#"),
		B.RandHex(6),
		B.TimeWithCode(DateFmtWeek, DateFmtMonth),
	)
	fnCnt := func(ctx context.Context, code string, min, step int64) (int64, error) {
		t.Log("累加器:", code, min, step)
		v, _ := map_cnt[code]
		if v < min {
			v = min
		} else {
			v += step
		}
		map_cnt[code] = v
		return v, nil
	}
	sg, err := NewGenerator("", exper, func(ctx context.Context, name string) (string, error) {
		return "{" + strings.ToUpper(name) + "}", nil
	}, fnCnt)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()

	for i := 1; i < 20; i++ {
		v, e := sg.Next(ctx)
		t.Log(v, e)
	}

	exper = B.Join("",
		B.CodeOnly(B.Const("order")),
		B.Code(B.Env(B.Const("STORE_ID"))),
		B.TimeWithCode(DateFmtMonth, ""),
		B.RandFillDigit(B.Incr(1, 1).SetAlign(6, 1, "@"), 10),
	)
	map_env := map[string]string{}
	sg, err = NewGenerator("", exper, NewMapEnv(map_env), fnCnt)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 1; i < 100; i++ {
		map_env["STORE_ID"] = fmt.Sprintf("%02d", i%4+1)
		v, e := sg.Next(ctx)
		t.Log(v, e)
	}
}
