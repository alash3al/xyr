package sql

import (
	"github.com/alash3al/xyr/internals/kernel"

	// the sql supported drivers
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/SAP/go-hdb/driver"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/sijms/go-ora/v2"
)

func init() {
	kernel.RegisterImporter("sqlite3", NewSQLDriver("sqlite3"))
	kernel.RegisterImporter("postgres", NewSQLDriver("postgres"))
	kernel.RegisterImporter("mysql", NewSQLDriver("mysql"))
	kernel.RegisterImporter("clickhouse", NewSQLDriver("clickhouse"))
	kernel.RegisterImporter("hana", NewSQLDriver("hdb"))
	kernel.RegisterImporter("oracle", NewSQLDriver("oracle"))
	kernel.RegisterImporter("sqlserver", NewSQLDriver("sqlserver"))
}
