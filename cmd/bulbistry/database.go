package main

import (
	"internal/config"
	"internal/database"

	"github.com/urfave/cli/v2"
)

func InitializeDatabase(ctx *cli.Context) error {
	cfg, err := config.ReadConfig(ctx.String("config"))
	if err != nil {
		return err
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		return err
	}

	if err = db.InitializeDatabase(); err != nil {
		return err
	}

	return nil
}

func MigrateDatabase(ctx *cli.Context) error {
	cfg, err := config.ReadConfig(ctx.String("config"))
	if err != nil {
		return err
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		return err
	}

	if err = db.MigrateDatabase(); err != nil {
		return err
	}

	return nil
}
