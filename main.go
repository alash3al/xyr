package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alash3al/xyr/config"
	"github.com/alash3al/xyr/driver"
	"github.com/alash3al/xyr/utils"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"

	_ "github.com/alash3al/xyr/driver/drivers/localjson"
	_ "github.com/mattn/go-sqlite3"
)

var (
	cfg    *config.Config
	dbConn *sqlx.DB

	tables = map[string]*config.Table{}
)

func main() {
	var err error

	cfg, err = config.LoadConfigFromFile(utils.Getenv("XYRCONFIG", "./config.xyr.hcl"))
	if err != nil {
		log.Fatal(err)
	}

	info, err := os.Stat(cfg.DataDir)
	if err != os.ErrNotExist {
		err = os.MkdirAll(cfg.DataDir, 0755)
	}

	if err != nil {
		log.Fatal(err)
	}

	if info != nil && !info.IsDir() {
		log.Fatal(fmt.Errorf("the specified path (%s) isn't a valid directory", cfg.DataDir))
	}

	dbConn, err = sqlx.Connect("sqlite3", filepath.Join(cfg.DataDir, "db.xyr")+"?_journal_mode=wal")
	if err != nil {
		log.Fatal(err)
	}

	for _, tb := range cfg.Tables {
		d, err := driver.Open(tb.DSN)
		if err != nil {
			log.Fatal(err)
		}

		if len(tb.Columns) < 1 {
			log.Fatal(fmt.Errorf("there is no columns for table %s", tb.Name))
		}

		tables[tb.Name] = tb
		tables[tb.Name].DriverInstance = d
	}

	app := &cli.App{
		Name:  "xyr",
		Usage: "query multiple data sources using sql",
		Commands: []*cli.Command{
			{
				Name:  "import",
				Usage: "import the defined tables data into xyr",
				Action: func(c *cli.Context) error {
					wg := &sync.WaitGroup{}

					for name, tb := range tables {
						if _, err := dbConn.Exec("DROP TABLE IF EXISTS " + name); err != nil {
							log.Fatal(err)
						}

						if _, err := dbConn.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", name, strings.Join(tb.Columns, ","))); err != nil {
							log.Fatal(err)
						}

						wg.Add(1)
						go (func(name string, tb *config.Table) {
							defer wg.Done()
							resultChan, errChan, doneChan := tb.DriverInstance.Run(tb.Loader)

						eventLoop:
							for {
								select {
								case <-doneChan:
									break eventLoop
								case err := <-errChan:
									log.Println(err)
								case result := <-resultChan:
									filteredResult := map[string]interface{}{}
									placeholders := []string{}
									for _, col := range tb.Columns {
										val, exists := result[col]
										if !exists {
											continue
										}
										filteredResult[col] = val
										placeholders = append(placeholders, ":"+col)
									}
									querySQL := fmt.Sprintf("INSERT INTO %s VALUES(%s)", name, strings.Join(placeholders, ","))
									if _, err := dbConn.NamedExec(querySQL, filteredResult); err != nil {
										log.Println(err)
										continue
									}
								}
							}

						})(name, tb)
					}

					wg.Wait()

					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}
