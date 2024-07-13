package jgorm

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
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

////////////////////////////////////////////////////////////////

type TableOptions interface {
	TableOptions() string
}

type TableComment interface {
	TableComment() string
}

var commentEscaper = strings.NewReplacer(
	`'`, "''",
	`\n`, "\\n",
	`\r`, "\\r",
)

// 处理表注释，只对mysql的新增表有效，其他数据库忽略
//
// todo: 暂时只处理新增时的表注释，后面再加上修改表注释
func AutoMigrate(db *gorm.DB, dst ...interface{}) (err error) {
	for _, v := range dst {
		var options string
		if v1, ok := v.(TableOptions); ok {
			options = v1.TableOptions()
		}
		if v1, ok := v.(TableComment); ok {
			if c1 := v1.TableComment(); c1 != "" {
				c1 = commentEscaper.Replace(c1)
				options = fmt.Sprintf(" COMMENT '%s'", c1)
			}
		}
		if options != "" {
			err = db.Set("gorm:table_options", options).AutoMigrate(v)
		} else {
			err = db.AutoMigrate(v)
		}
		if err != nil {
			return
		}
	}
	return nil
}

// FindInBatches4Ordered 从已经排序的查询中批量取数据, gorm.DB#FindInBatches 需要有唯一主键
func FindInBatches4Ordered(db1 *gorm.DB, dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) *gorm.DB {
	var (
		tx           = db1.Session(&gorm.Session{})
		queryDB      = tx
		rowsAffected int64
		batch        int
	)

	// user specified offset or limit
	var totalSize int
	if c, ok := tx.Statement.Clauses["LIMIT"]; ok {
		if limit, ok := c.Expression.(clause.Limit); ok {
			if limit.Limit != nil {
				totalSize = *limit.Limit
			}

			if totalSize > 0 && batchSize > totalSize {
				batchSize = totalSize
			}

			// reset to offset to 0 in next batch
			tx = tx.Offset(-1).Session(&gorm.Session{})
		}
	}

	for {
		result := queryDB.Offset(int(rowsAffected)).Limit(batchSize).Find(dest)
		rowsAffected += result.RowsAffected
		batch++

		if result.Error == nil && result.RowsAffected != 0 {
			fcTx := result.Session(&gorm.Session{NewDB: true})
			fcTx.RowsAffected = result.RowsAffected
			_ = tx.AddError(fc(fcTx, batch))
		} else if result.Error != nil {
			_ = tx.AddError(result.Error)
		}

		if tx.Error != nil || int(result.RowsAffected) < batchSize {
			break
		}

		if totalSize > 0 {
			if totalSize <= int(rowsAffected) {
				break
			}
			if totalSize/batchSize == batch {
				batchSize = totalSize % batchSize
			}
		}
	}

	tx.RowsAffected = rowsAffected
	return tx
}
