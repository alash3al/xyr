package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/alash3al/xyr/internals/kernel"
	"github.com/alash3al/xyr/utils"
	"github.com/urfave/cli/v2"
)

func init() {
	RegisterCommand(Exec)
}

// Exec executes an sql statement
func Exec(env *kernel.Env) *cli.Command {
	return &cli.Command{
		Name:      "exec",
		Usage:     "execute the specified sql statement",
		UsageText: "xyr exec 'SELECT * FROM MY_TABLE'",
		Action: func(c *cli.Context) error {
			stmnt := strings.Join(c.Args().Slice(), " ")
			if strings.TrimSpace(stmnt) == "" {
				return fmt.Errorf("empty sql statment specified")
			}

			log.Printf("executing %s", stmnt)

			rows, err := env.DBConn.Queryx(stmnt)
			if err != nil {
				return err
			}

			result := utils.SqlAll(rows)
			jsonResult, _ := json.MarshalIndent(result, "", " ")

			fmt.Println(string(jsonResult))

			return nil
		},
	}
}
