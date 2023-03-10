package main

import (
	"log"

	"github.com/ceit-aut/S7IE03/cmd"
	"github.com/ceit-aut/S7IE03/internal/config"
	"github.com/ceit-aut/S7IE03/internal/repositories"
	"github.com/ceit-aut/S7IE03/internal/storage"
	"github.com/ceit-aut/S7IE03/pkg/auth"

	"github.com/spf13/cobra"
)

func main() {
	// create main command
	root := cobra.Command{}

	// load configs
	cfg := config.Load()

	// open connection to mongodb
	stg, err := storage.NewConnection(cfg.Storage)
	if err != nil {
		log.Printf("connection to database failed:\n\t%v\n", err)
	}

	// creating auth client
	au := auth.New(cfg.JWT)

	// create repositories
	rp := repositories.New(stg)

	// add cobra commands
	root.AddCommand(
		cmd.HTTP{Cfg: cfg, Repositories: *rp, Auth: au}.Command(),
		cmd.Worker{Cfg: cfg.Worker, Repositories: *rp}.Command(),
	)

	// execute command
	if er := root.Execute(); er != nil {
		panic(er)
	}
}
