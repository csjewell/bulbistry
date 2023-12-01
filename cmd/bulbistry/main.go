// The command bulbistry implements a minimal compliant OCI container registry
// suitable for self-hosting in containers.
package main

import (
	"internal/config"
	"internal/database"
	v "internal/version"

	"fmt"
	"os"

	tbv "github.com/csjewell/bulbistry"
	htpasswd "github.com/tg123/go-htpasswd"
	"github.com/urfave/cli/v2"
)

var authorizer *htpasswd.File
var executionLog *tbv.Logger
var debugLog *tbv.DebugLogger
var cfg config.Config
var db database.Database

func main() {
	executionLog = tbv.NewLogger(os.Stderr, nil)
	debugLog = tbv.NewDebugLogger(executionLog)
	debugLog.Print("Bulbistry started")

	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Printf("%s\n", ctx.App.Version)
	}

	app := &cli.App{
		Name:      "Bulbistry",
		HelpName:  "bulbistry",
		Usage:     "A pared-down container registry, perfect for self-hosting",
		Authors:   []*cli.Author{{Name: "Curtis Jewell", Email: "swordsman@curtisjewell.name"}},
		Copyright: "Copyright (c) 2023 Curtis Jewell",
		Version:   v.Version(),
		Action:    func(c *cli.Context) error { return RunServer(c) },
		Commands: []*cli.Command{
			&cli.Command{
				Name:      "create-config",
				Aliases:   []string{"config"},
				Usage:     "Generate a basic configuration",
				UsageText: "Generates a configuration interactively",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "config",
						Aliases:   []string{"c"},
						Value:     "bulbistry.yml",
						Usage:     "Save configuration to FILE",
						TakesFile: true,
						EnvVars:   []string{"BULBISTRY_CONFIG"},
					},
				},
				Action: func(c *cli.Context) error {
					return CreateConfig(c)
				},
			},
			&cli.Command{
				Name:      "init-db",
				Aliases:   []string{"init-database"},
				Usage:     "Initialize application database",
				UsageText: "Initialize application database",
				Action: func(c *cli.Context) error {
					return InitializeDatabase(c)
				},
			},
			&cli.Command{
				Name:      "migrate-db",
				Aliases:   []string{"init-database"},
				Usage:     "Migrate the database",
				UsageText: "Migrate the database into an upgraded format",
				Action: func(c *cli.Context) error {
					return MigrateDatabase(c)
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Value:   false,
				Usage:   "Turns on debug logging",
				Action: func(c *cli.Context, b bool) error {
					if b {
						return debugLog.TurnOn()
					} else {
						return debugLog.TurnOff()
					}
				},
			},
			&cli.StringFlag{
				Name:      "config",
				Aliases:   []string{"c"},
				Value:     "bulbistry.yml",
				Usage:     "Load configuration from FILE",
				TakesFile: true,
				EnvVars:   []string{"BULBISTRY_CONFIG"},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		executionLog.Output(2, err.Error())
		os.Exit(1)
	}
}
