package kernel

import "github.com/jmoiron/sqlx"

type Env struct {
	Config *Config
	Tables map[string]*Table
	DBConn *sqlx.DB
}
