package sql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Driver represents the main importer driver
type Driver struct {
	sqlDriverName string
	conn          *sqlx.DB
}

func NewSQLDriver(driverName string) *Driver {
	return &Driver{
		sqlDriverName: driverName,
	}
}

// Open implements Importer#open
func (d *Driver) Open(dsn string) error {
	if d.sqlDriverName == "" {
		return fmt.Errorf("unable to detect the required sql driver, this means you didn't make use of NewSQLDriver(...)")
	}

	conn, err := sqlx.Connect(d.sqlDriverName, dsn)
	if err != nil {
		return err
	}

	d.conn = conn

	return nil
}

// Import implements Importer#import
func (d *Driver) Import(sqlStmnt string) (<-chan map[string]interface{}, <-chan error, <-chan bool) {
	resultChan := make(chan map[string]interface{})
	errChan := make(chan error)
	doneChan := make(chan bool)

	go (func() {
		defer (func() {
			doneChan <- true

			close(resultChan)
			close(errChan)
			close(doneChan)
		})()

		rows, err := d.conn.Queryx(sqlStmnt)
		if err != nil {
			errChan <- err
			return
		}

		for rows.Next() {
			row := map[string]interface{}{}

			if err := rows.MapScan(row); err != nil {
				fmt.Println(err)
				errChan <- err
				continue
			}

			resultChan <- row
		}
	})()

	return resultChan, errChan, doneChan
}
