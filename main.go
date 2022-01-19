package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/alash3al/xyr/internals/commands"
	"github.com/alash3al/xyr/internals/kernel"
	"github.com/alash3al/xyr/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/dnscache"
	"github.com/urfave/cli/v2"

	_ "github.com/alash3al/xyr/internals/importers/jsondir"
	_ "github.com/alash3al/xyr/internals/importers/s3jsondir"
	_ "github.com/alash3al/xyr/internals/importers/sql"
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
			kernel.Logger.Fatal(err)
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
			kernel.Logger.Fatal(err)
		}

		if info != nil && !info.IsDir() {
			kernel.Logger.Fatal(fmt.Errorf("the specified path (%s) isn't a valid directory", kernelEnv.Config.DataDir))
		}
	}

	// workers count
	{
		if kernelEnv.Config.WorkersCount < 1 {
			kernelEnv.Config.WorkersCount = runtime.NumCPU()
		}
	}

	// logger
	{
		if !kernelEnv.Config.Debug {
			kernel.Logger.WithoutDebug()
		}
	}

	// Initialize the storage engine
	{
		dbConn, err := sqlx.Connect("sqlite3", filepath.Join(kernelEnv.Config.DataDir, "db.xyr")+"?_journal_mode=wal&cache=shared&_sync=0")
		if err != nil {
			kernel.Logger.Fatal(err)
		}

		kernelEnv.DBConn = dbConn
	}

	// cache the tables inside our kernel environment container
	{
		for _, tb := range kernelEnv.Config.Tables {
			d, err := kernel.OpenImporter(tb.ImporterName, tb.DSN)
			if err != nil {
				kernel.Logger.Fatal(err)
			}

			if len(tb.Columns) < 1 {
				kernel.Logger.Fatal(fmt.Errorf("there is no columns for table %s", tb.Name))
			}

			kernelEnv.Tables[tb.Name] = tb
			kernelEnv.Tables[tb.Name].ImporterInstance = d
		}
	}

	// default http transport dns caching
	{
		r := &dnscache.Resolver{}
		http.DefaultClient.Transport = &http.Transport{
			MaxIdleConnsPerHost: 64,
			DialContext: func(ctx context.Context, network string, addr string) (conn net.Conn, err error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}
				ips, err := r.LookupHost(ctx, host)
				if err != nil {
					return nil, err
				}
				for _, ip := range ips {
					var dialer net.Dialer
					conn, err = dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
					if err == nil {
						break
					}
				}
				return
			},
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
