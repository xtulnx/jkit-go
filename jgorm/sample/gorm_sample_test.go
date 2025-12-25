package sample

import (
	"database/sql"
	"testing"

	"github.com/xtulnx/jkit-go/jgorm"
	"github.com/xtulnx/jkit-go/jgorm/config"
	"gorm.io/datatypes"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type Foo1 struct {
	gorm.Model
	StoreId     uint           `gorm:"type:int;index:s;index:sd,unique;comment:店铺 id"`
	BusinessDay sql.NullTime   `gorm:"type:date;comment:营业日期"`
	Ext         datatypes.JSON `gorm:"type:json;comment:扩展信息"`
	StoreName   string         `gorm:"type:varchar(128);comment:店铺名称"`
}
type Foo1Fields struct {
	jgorm.SimpleFields
	StoreId     field.Uint //店铺id
	BusinessDay field.Time //营业日
	Ext         field.Field
	StoreName   field.String
}

func (i *Foo1Fields) UpdateTableName(table string) {
	i.SimpleFields.UpdateTableName(table)
	i.StoreId = field.NewUint(table, "store_id")
	i.BusinessDay = field.NewTime(table, "business_day")
	i.Ext = field.NewField(table, "ext")
	i.StoreName = field.NewString(table, "store_name")
}

func TestDaoWrapper(t *testing.T) {
	db0, err := OpenDb(config.NewDbProvider1("sqlite", "file::memory:?cache=shared"))
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range []struct {
		N    string
		F    func() string
		Want string
	}{
		{
			"dao 的默认查询 sql", func() string {
				db1 := db0.Model(&Foo1{})
				return jgorm.ShowSQL(db1)
			}, "SELECT * FROM `foo1` WHERE `foo1`.`deleted_at` IS NULL",
		},
		{
			"dao 的默认查询 sql，去除 scope", func() string {
				db1 := db0.Unscoped().Model(&Foo1{})
				return jgorm.ShowSQL(db1)
			}, "SELECT * FROM `foo1`",
		},
		{
			"dao的复杂查询，自定义 where 与 select", func() string {
				dao1 := jgorm.NewDaoWrapper(db0.Model(&Foo1{}), nil, "t")
				f1 := &Foo1Fields{}
				f1.UpdateTableName("t")
				dao1_1 := dao1.GetDaoDirect().Select(f1.StoreId, f1.StoreName)
				// 补充 where 条件
				jgorm.DbClause(dao1_1,
					f1.StoreName.Like("A%"), jgorm.DbClauseExpr("?->>'$.meta.income > ?", f1.Ext, 1000),
				)
				// 增加自定义检索字段
				jgorm.SelectAppend(dao1_1.(*gen.DO),
					jgorm.ClauseExpr("?->>'$.meta.income' as income", f1.Ext),                                      // 提取 json 字段
					jgorm.ClauseExpr("concat_ws(':',date_format(?,'%Y-%m-%d'),?) k", f1.BusinessDay, f1.StoreName), // 构造复杂字段
				)
				return jgorm.DaoShowSQL(dao1_1)
			}, "SELECT `t`.`store_id`,`t`.`store_name`,`t`.`ext`->>'$.meta.income' as income,concat_ws(':',date_format(`t`.`business_day`,'%Y-%m-%d'),`t`.`store_name`) k FROM (SELECT * FROM `foo1` WHERE `foo1`.`deleted_at` IS NULL) t WHERE `t`.`store_name` LIKE \"A%\" AND `t`.`ext`->>'$.meta.income > 1000",
		},
	} {
		v1 := v
		t.Run(v.N, func(t *testing.T) {
			s1 := v1.F()
			if s1 != v1.Want {
				t.Errorf("%s => [%s], want [%s]", v1.N, s1, v1.Want)
			}
		})
	}
}
