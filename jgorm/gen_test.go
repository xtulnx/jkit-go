package jgorm

import (
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
	"testing"
	"unsafe"
	_ "unsafe"
)

//go:linkname setE gorm.io/gen/field/(expr).setE
func setE(e expr, v clause.Expression) expr

type expr struct {
	col       clause.Column
	e         clause.Expression
	buildOpts []field.BuildOpt
}

func TestGenExpr(t *testing.T) {
	e1 := field.EmptyExpr()
	t.Log(e1)
	e2 := *(*expr)(unsafe.Pointer(&e1))
	e3 := setE(e2, clause.Eq{Column: "a", Value: 1})
	//e1 = setE(expr(e1), clause.Eq{Column: "a", Value: 1})
	t.Log(e3)
}
