package jgorm

import (
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// gorm/gen 的查询封装，主要是别名处理
//
// jason.liao

type DaoWrapper interface {
	GetAlias() string
	SetAlias(alias string)
	GetDbRelated() *gorm.DB
	GetDaoRelated() gen.Dao
	GetDbDirect() *gorm.DB
	GetDaoDirect() gen.Dao
}

type tDaoWrapper struct {
	db0   *gorm.DB
	db1   *gorm.DB
	alias string

	// 用于联表
	dbRelated  *gorm.DB
	daoRelated *gen.DO
	// 用于直接查询
	dbDirect  *gorm.DB
	daoDirect *gen.DO
}

func NewDaoWrapper1(db1 *gorm.DB, alias string) DaoWrapper {
	return NewDaoWrapper(nil, db1, alias)
}

// NewDaoWrapperByJoin 通过联表查询构建 DaoWrapper
//
//   - dao1 已经是可以用于联表查询的 daoRelated
func NewDaoWrapperByJoin(dao1 gen.Dao, alias string) DaoWrapper {
	return NewDaoWrapper(nil, dao1.(*gen.DO).UnderlyingDB(), alias)
}

func NewDaoWrapper(db0, db1 *gorm.DB, alias string) DaoWrapper {
	if db0 != nil && db1 == nil {
		db1, db0 = db0, db1
	}
	if db0 == nil {
		db0 = db1.Session(&gorm.Session{Initialized: true, NewDB: true})
	}
	d1 := &tDaoWrapper{
		db0:        db0,
		db1:        db1,
		alias:      alias,
		dbRelated:  nil,
		daoRelated: nil,
		dbDirect:   nil,
		daoDirect:  nil,
	}
	return d1
}

func (D *tDaoWrapper) GetAlias() string {
	return D.alias
}

// SetAlias 修改别名
func (D *tDaoWrapper) SetAlias(alias string) {
	if D.alias != alias {
		D.alias = alias
		D.dbDirect = nil
		D.daoDirect = nil
		D.dbRelated = nil
		D.daoRelated = nil
	}
}

// GetDbRelated 用于联表查询
func (D *tDaoWrapper) GetDbRelated() *gorm.DB {
	if D.dbRelated == nil {
		D.dbRelated = D.db0.Table("?", D.db1)
		D.dbRelated.Statement.WriteQuoted(clause.Table{Name: clause.CurrentTable})
	}
	return D.dbRelated
}

// GetDaoRelated 用于联表查询
func (D *tDaoWrapper) GetDaoRelated() gen.Dao {
	if D.daoRelated == nil {
		dbR := D.GetDbRelated()
		daoR := &gen.DO{}
		daoR.UseDB(dbR)
		D.daoRelated = daoR.As(D.alias).(*gen.DO)
	}
	return D.daoRelated
}

// GetDbDirect 用于直接查询
func (D *tDaoWrapper) GetDbDirect() *gorm.DB {
	if D.dbDirect == nil {
		D.dbDirect = D.db0.Table("(?) "+D.alias, D.db1)
	}
	return D.dbDirect
}

// GetDaoDirect 用于直接查询，不能使用 gen.Table
func (D *tDaoWrapper) GetDaoDirect() gen.Dao {
	if D.daoDirect == nil {
		dbD := D.GetDbDirect()
		daoD := &gen.DO{}
		daoD.UseDB(dbD)
		D.daoDirect = daoD
	}
	return D.daoDirect
}
