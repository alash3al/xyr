package utils

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func SqlAll(rows *sqlx.Rows) []map[string]interface{} {
	result := []map[string]interface{}{}
	for rows.Next() {
		m := map[string]interface{}{}

		if err := rows.MapScan(m); err != nil {
			log.Println(err)
		}

		result = append(result, m)
	}
	return result
}
