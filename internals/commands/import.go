package commands

import (
	"encoding/json"
	"fmt"
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
		Name:      "table:import",
		Usage:     "import the defined tables data into xyr",
		UsageText: "xyr table:import table1 table2 ...",
		Action: func(c *cli.Context) error {
			wg := &sync.WaitGroup{}
			selectedTables := c.Args().Slice()

			if len(selectedTables) < 1 {
				return fmt.Errorf("no table specified")
			}

			for _, selectedTable := range selectedTables {
				tb, found := env.Tables[selectedTable]

				if !found {
					kernel.Logger.Fatal(fmt.Sprintf("undefined table %s", selectedTable))
				}

				if _, err := env.DBConn.Exec("DROP TABLE IF EXISTS " + tb.Name); err != nil {
					kernel.Logger.Fatal(err)
				}

				if _, err := env.DBConn.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tb.Name, strings.Join(tb.Columns, ","))); err != nil {
					kernel.Logger.Fatal("unable to create table", err)
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
	resultChan, errChan, doneChan := tb.ImporterInstance.Import(tb.Filter)
	loop := true

	for loop {
		select {
		case <-doneChan:
			loop = false
			break
		case err := <-errChan:
			kernel.Logger.Error(err.Error())
		case result := <-resultChan:
			kernel.Logger.Info("processing", result)

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
				kernel.Logger.Error("unable to find a document to be written")
				continue
			}
			querySQL := fmt.Sprintf("INSERT INTO %s VALUES(%s)", tb.Name, strings.Join(placeholders, ","))
			if _, err := env.DBConn.NamedExec(querySQL, filteredResult); err != nil {
				kernel.Logger.Error(querySQL, err)
				continue
			}
		}
	}

	kernel.Logger.Info(fmt.Sprintf("Congrats! now you can execute `xyr exec 'SELECT * FROM %s'`", tb.Name))
}
