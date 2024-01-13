package jgorm

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

// 补充 gorm/gen 辅助

// 合并多个查询表达式，转成 sql 与 参数对象
//
//	columns: field.Expr, string, clause.Expr
func buildExpr(stmt *gorm.Statement, columns ...interface{}) (query []string, args []interface{}) {
	for _, e := range columns {
		switch v := e.(type) {
		case field.Expr:
			sql, vars := v.BuildWithArgs(stmt)
			query = append(query, sql.String())
			args = append(args, vars...)
		case string:
			query = append(query, v)
		case clause.Expr:
			query = append(query, v.SQL)
			args = append(args, v.Vars...)
		}
	}
	return query, args
}

// Select 字段
//
//	columns: field.Expr, string, clause.Expr
func Select(do *gen.DO, columns ...interface{}) {
	db := do.UnderlyingDB()
	query, args := buildExpr(db.Statement, columns...)
	db = db.Select(strings.Join(query, ","), args...)
	do.ReplaceDB(db)
}

// SelectAppend 增加字段
//
//	columns: field.Expr, string, clause.Expr
func SelectAppend(do *gen.DO, columns ...interface{}) {
	db := do.UnderlyingDB()
	query, args := buildExpr(db.Statement, columns...)
	if c1, ok := db.Statement.Clauses["SELECT"]; ok && c1.Expression != nil {
		switch v := c1.Expression.(type) {
		case clause.Expr:
			query = append([]string{v.SQL}, query...)
			args = append(v.Vars, args...)
		case clause.NamedExpr:
			query = append([]string{v.SQL}, query...)
			args = append(v.Vars, args...)
		}
	} else {
		query = append(db.Statement.Selects, query...)
	}
	if do.TableName() != "" && db.Statement.TableExpr == nil {
		db = db.Table(do.TableName())
	}
	db = db.Select(strings.Join(query, ","), args...)
	do.ReplaceDB(db)
}

func ColsNamesByExpr(expr ...field.Expr) []string {
	names := make([]string, len(expr))
	for i := range expr {
		names[i] = string(expr[i].ColumnName())
	}
	return names
}

// 使用别名构建子查询
func DaoToSub(dao gen.Dao, alias string) gen.Dao {
	db1 := dao.(*gen.DO).UnderlyingDB()
	db0 := db1.Session(&gorm.Session{Initialized: true, NewDB: true})
	db1 = db0.Table("(?) "+alias, db1)
	d2 := &gen.DO{}
	d2.UseDB(db1)
	return d2
}

// 获取一个空的 db，用于构建子查询
func DaoDbBlank(dao gen.Dao) *gorm.DB {
	return dao.(*gen.DO).UnderlyingDB().Session(&gorm.Session{Initialized: true, NewDB: true})
}

// 获取 dao 的 sql
func DaoShowSQL(dao gen.Dao) string {
	return dao.(*gen.DO).UnderlyingDB().ToSQL(func(tx *gorm.DB) *gorm.DB {
		var tm interface{}
		return tx.Scan(&tm)
	})
}
