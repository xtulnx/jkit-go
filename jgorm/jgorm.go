package jgorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
)

func StmtReplaceColumnValue(stmt *gorm.Statement, field *schema.Field, fn func(r1 interface{}, zero bool) (r2 interface{}, replace bool)) {
	if v, ok := stmt.Dest.(map[string]interface{}); ok {
		r1, ok1 := v[field.Name]
		if r2, replace := fn(r1, !ok1); replace {
			v[field.Name] = r2
		}
		return
	}

	if v, ok := stmt.Dest.([]map[string]interface{}); ok {
		for _, m := range v {
			r1, ok1 := m[field.Name]
			if r2, replace := fn(r1, !ok1); replace {
				m[field.Name] = r2
			}
		}
		return
	}

	if stmt.Schema != nil {
		destValue := reflect.ValueOf(stmt.Dest)
		for destValue.Kind() == reflect.Ptr {
			destValue = destValue.Elem()
		}

		if stmt.ReflectValue != destValue {
			if !destValue.CanAddr() {
				destValueCanAddr := reflect.New(destValue.Type())
				destValueCanAddr.Elem().Set(destValue)
				stmt.Dest = destValueCanAddr.Interface()
				destValue = destValueCanAddr.Elem()
			}

			switch destValue.Kind() {
			case reflect.Struct:
				r1, zero := field.ValueOf(stmt.Context, destValue)
				if r2, replace := fn(r1, zero); replace {
					stmt.AddError(field.Set(stmt.Context, destValue, r2))
				}
			default:
				stmt.AddError(gorm.ErrInvalidData)
			}
		}

		switch stmt.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < stmt.ReflectValue.Len(); i++ {
				r1, zero := field.ValueOf(stmt.Context, stmt.ReflectValue.Index(i))
				if r2, replace := fn(r1, zero); replace {
					stmt.AddError(field.Set(stmt.Context, stmt.ReflectValue.Index(i), r2))
				}
			}
		case reflect.Struct:
			if !stmt.ReflectValue.CanAddr() {
				stmt.AddError(gorm.ErrInvalidData)
				return
			}
			r1, zero := field.ValueOf(stmt.Context, stmt.ReflectValue)
			if r2, replace := fn(r1, zero); replace {
				stmt.AddError(field.Set(stmt.Context, stmt.ReflectValue, r2))
			}
		}
	} else {
		stmt.AddError(gorm.ErrInvalidData)
	}
}
