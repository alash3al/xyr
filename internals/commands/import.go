package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/alash3al/xyr/internals/kernel"
	"github.com/urfave/cli/v2"
)

func init() {
	RegisterCommand(Import)
}

// Import implements the import sub command
func Import(env *kernel.Env) *cli.Command {
	return &cli.Command{
		Name:  "table:import",
		Usage: "import the defined tables data into xyr",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "table",
				Aliases:  []string{"t"},
				Usage:    "specify which table(s) to be loaded, this flag could be specified multiple times",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			wg := &sync.WaitGroup{}
			selectedTables := c.StringSlice("table")

			for _, selectedTable := range selectedTables {
				tb, found := env.Tables[selectedTable]

				if !found {
					log.Fatal(fmt.Sprintf("undefined table %s", tb.Name))
				}

				if _, err := env.DBConn.Exec("DROP TABLE IF EXISTS " + tb.Name); err != nil {
					log.Fatal(err)
				}

				if _, err := env.DBConn.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tb.Name, strings.Join(tb.Columns, ","))); err != nil {
					log.Fatal("unable to create table", err)
				}

				wg.Add(1)

				go (func(tb *kernel.Table) {
					defer wg.Done()

					load(tb, env)
				})(tb)
			}

			wg.Wait()

			return nil
		},
	}
}

func load(tb *kernel.Table, env *kernel.Env) {
	resultChan, errChan, doneChan := tb.ImporterInstance.Import(tb.Loader)
	loop := true

	for loop {
		select {
		case <-doneChan:
			loop = false
			break
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
				switch val.(type) {
				case []interface{}, map[string]interface{}:
					val, _ = json.Marshal(val)
				}
				filteredResult[col] = val
				placeholders = append(placeholders, ":"+col)
			}
			if len(filteredResult) < 1 {
				fmt.Println(result)
				log.Println("unable to find a document to be written")
				continue
			}
			querySQL := fmt.Sprintf("INSERT INTO %s VALUES(%s)", tb.Name, strings.Join(placeholders, ","))
			if _, err := env.DBConn.NamedExec(querySQL, filteredResult); err != nil {
				log.Println(querySQL, err)
				continue
			}
		}
	}
}
