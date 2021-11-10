package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alash3al/xyr/internals/commands"
	"github.com/alash3al/xyr/internals/kernel"
	"github.com/alash3al/xyr/utils"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"

	_ "github.com/alash3al/xyr/internals/importers/localjsonobj"
	_ "github.com/mattn/go-sqlite3"
)

var (
	kernelEnv = &kernel.Env{
		Tables: map[string]*kernel.Table{},
	}
)

func main() {
	// load config file and parse it
	{
		cfg, err := kernel.LoadConfigFromFile(utils.Getenv("XYRCONFIG", "./config.xyr.hcl"))
		if err != nil {
			log.Fatal(err)
		}

		kernelEnv.Config = cfg
	}

	// creates the main data directory if not exists
	{
		info, err := os.Stat(kernelEnv.Config.DataDir)
		if err != os.ErrNotExist {
			err = os.MkdirAll(kernelEnv.Config.DataDir, 0755)
		}

		if err != nil {
			log.Fatal(err)
		}

		if info != nil && !info.IsDir() {
			log.Fatal(fmt.Errorf("the specified path (%s) isn't a valid directory", kernelEnv.Config.DataDir))
		}
	}

	// Initialize the storage engine
	{
		dbConn, err := sqlx.Connect("sqlite3", filepath.Join(kernelEnv.Config.DataDir, "db.xyr")+"?_journal_mode=wal")
		if err != nil {
			log.Fatal(err)
		}

		kernelEnv.DBConn = dbConn
	}

	// cache the tables inside our kernel environment container
	{
		for _, tb := range kernelEnv.Config.Tables {
			d, err := kernel.OpenImporter(tb.DSN)
			if err != nil {
				log.Fatal(err)
			}

			if len(tb.Columns) < 1 {
				log.Fatal(fmt.Errorf("there is no columns for table %s", tb.Name))
			}

			kernelEnv.Tables[tb.Name] = tb
			kernelEnv.Tables[tb.Name].ImporterInstance = d
		}
	}

	app := &cli.App{
		Name:     "xyr",
		Usage:    "query multiple data sources using sql",
		Commands: []*cli.Command{},
	}

	for _, cmdFactory := range commands.GetRegisteredCommands() {
		app.Commands = append(app.Commands, cmdFactory(kernelEnv))
	}

	app.Run(os.Args)
}
