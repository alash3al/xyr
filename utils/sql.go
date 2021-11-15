package utils

import (
	"github.com/alash3al/xyr/internals/kernel"
	"github.com/jmoiron/sqlx"
)

func SqlAll(rows *sqlx.Rows) []map[string]interface{} {
	result := []map[string]interface{}{}
	for rows.Next() {
		m := map[string]interface{}{}

		if err := rows.MapScan(m); err != nil {
			kernel.Logger.Error(err)
		}

		result = append(result, m)
	}
	return result
}
