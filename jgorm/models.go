package jgorm

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"time"
)

////////////////////////////////////////////////////////////////

type SimpleDayType datatypes.Date

func kitDate2Str(t1 time.Time) string {
	if !t1.IsZero() {
		if t1.Location() != time.Local {
			t1 = t1.Local()
		}
		return t1.Format("2006-01-02")
	}
	return ""
}
func kitStr2Date(st string) time.Time {
	if st == "" {
		return time.Time{}
	}
	t1, err := time.ParseInLocation("2006-01-02", st, time.Local)
	if err != nil {
		return time.Time{}
	}
	return t1
}

func (dt SimpleDayType) String() string {
	return kitDate2Str(time.Time(dt))
}

func (d SimpleDayType) IsValid() bool {
	return !time.Time(d).IsZero()
}

func (dt *SimpleDayType) Scan(value interface{}) (err error) {
	var a sql.NullTime
	err = a.Scan(value)
	if err == nil && a.Valid {
		*dt = SimpleDayType(a.Time)
	}
	return
}

func (dt SimpleDayType) Value() (driver.Value, error) {
	t1 := time.Time(dt)
	if !t1.IsZero() {
		if t1.Location() != time.Local {
			t1 = t1.Local()
		}
		y, m, d := t1.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, t1.Location()), nil
	}
	return nil, nil
}

// GormDataType gorm common data type
func (dt SimpleDayType) GormDataType() string {
	return "date"
}

func (dt SimpleDayType) GobEncode() ([]byte, error) {
	return time.Time(dt).GobEncode()
}

func (dt *SimpleDayType) GobDecode(b []byte) error {
	return (*time.Time)(dt).GobDecode(b)
}

func (dt SimpleDayType) MarshalJSON() ([]byte, error) {
	s1 := dt.String()
	if s1 == "" {
		return []byte("null"), nil
	}
	return bytes.Join([][]byte{{'"'}, []byte(s1), {'"'}}, nil), nil
}

func (dt *SimpleDayType) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		return nil
	}
	b = bytes.TrimFunc(b, func(r rune) bool {
		return r == '"'
	})
	t1 := kitStr2Date(string(b))
	*dt = SimpleDayType(t1)
	return nil
}
func (dt SimpleDayType) MarshalText() ([]byte, error) {
	s1 := dt.String()
	return []byte(s1), nil
}
func (dt *SimpleDayType) UnmarshalText(data []byte) error {
	t1 := kitStr2Date(string(data))
	*dt = SimpleDayType(t1)
	return nil
}

// ---

// 特殊日期类型，在 db 执行 create 时自动计算一个合适的日期，如营业日期等
//
// jason.liao

func DefaultSimpleDay(curTime time.Time, colName, tableName string) time.Time {
	if y, m, d := curTime.Date(); curTime.Hour() < 9 {
		curTime = time.Date(y, m, d-1, 0, 0, 0, 0, curTime.Location())
	} else {
		curTime = time.Date(y, m, d, 0, 0, 0, 0, curTime.Location())
	}
	return curTime
}

// BusinessDayGenFn 自定义日期生成函数，参数：当前时间，字段名，表名
type BusinessDayGenFn func(curTime time.Time, colName, tableName string) time.Time

var defaultBusinessDayGenFn BusinessDayGenFn = DefaultSimpleDay

// SetDefaultBusinessDay 自定义业务日期生成函数
func SetDefaultBusinessDay(fn BusinessDayGenFn) {
	defaultBusinessDayGenFn = fn
}

func (n SimpleDayType) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SimpleDayCreateClause{Field: f}}
}

type SimpleDayCreateClause struct {
	Field *schema.Field
}

func (b SimpleDayCreateClause) Name() string {
	return ""
}

func (b SimpleDayCreateClause) Build(builder clause.Builder) {
}

func (b SimpleDayCreateClause) MergeClause(c *clause.Clause) {
}

func (b SimpleDayCreateClause) ModifyStatement(stmt *gorm.Statement) {
	var curTime time.Time
	StmtReplaceColumnValue(stmt, b.Field, func(r1 interface{}, zero bool) (r2 interface{}, replace bool) {
		if t1, ok := r1.(SimpleDayType); ok && t1.IsValid() {
			return
		}
		if curTime.IsZero() {
			curTime = stmt.DB.NowFunc()
			if fn := defaultBusinessDayGenFn; fn != nil {
				curTime = fn(curTime, b.Field.Name, b.Field.DBName)
			} else {
				curTime = DefaultSimpleDay(curTime, b.Field.Name, b.Field.DBName)
			}
		}
		return curTime, !curTime.IsZero()
	})
	if !curTime.IsZero() {
		stmt.AddClause(clause.Set{{Column: clause.Column{Name: b.Field.DBName}, Value: curTime}})
	}
}
