package config

import (
	"fmt"
	"testing"
)

func Test配置(t *testing.T) {
	for _, v := range []struct {
		N    string
		P    DsnProvider
		Want string
	}{
		{
			"mysql", &Mysql{
				GeneralDB: GeneralDB{
					FullDsn: "db-user:password@tcp(db-host:8904)/db-name?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
				},
			}, "[mysql]db-user:password@tcp(db-host:8904)/db-name?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
		},
		{
			"mysql2", &Mysql{
				GeneralDB: GeneralDB{
					Username: "db-user",
					Password: "password",
					Path:     "db-host",
					Port:     "8904",
					Dbname:   "db-name",
					Config:   "charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
				},
			}, "[mysql]db-user:password@tcp(db-host:8904)/db-name?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
		},
		{
			"mysql3", &SpecializedDB{
				GeneralDB: GeneralDB{
					Username: "db-user",
					Password: "password",
					Path:     "db-host",
					Port:     "8904",
					Dbname:   "db-name",
					Config:   "charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
				},
				Type: "mysql",
			}, "[mysql]db-user:password@tcp(db-host:8904)/db-name?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
		},
		{
			"postgres", &Pgsql{
				GeneralDB: GeneralDB{
					Path:     "127.0.0.1",
					Port:     "5432",
					Dbname:   "mydb",
					Username: "user",
					Password: "pwd",
					Config:   "sslmode=disable",
				},
			}, "[postgres]host=127.0.0.1 user=user password=pwd dbname=mydb port=5432 sslmode=disable",
		},
		{
			"oracle", &SpecializedDB{
				GeneralDB: GeneralDB{
					Path:     "example.com",
					Port:     "1521",
					Dbname:   "XE",
					Username: "scott",
					Password: "tiger",
					Config:   "",
				},
				Type: "oracle",
			}, "[oracle]oracle://scott:tiger@example.com:1521/XE?",
		},
	} {
		v1 := v
		t.Run(v.N, func(t *testing.T) {
			g1, ok := ExtractGeneralDB(v1.P)
			if !ok {
				t.Errorf("%s => %#v, invalid GeneralDB", v1.N, v1.P)
			} else if g1.LogZap != "" {
				//
			}
			s1 := fmt.Sprintf("[%s]%s", v1.P.Dialect(), v1.P.Dsn())
			if s1 != v1.Want {
				t.Errorf("%s => %#v, [%s], want [%s]", v1.N, v1.P, s1, v1.Want)
			}
		})
	}
}
