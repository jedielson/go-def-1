package main

import (
	"log"
	"os"

	"github.com/jedielson/bookstore/cmd/worker/actions"
	"github.com/jedielson/bookstore/cmd/worker/flags"
	"github.com/urfave/cli/v2"
)

var (
	AppName    = "go-bookstore"
	AppUsage   = ""
	AppVersion = "0.0.1"
	// GitSummary = "none"
	// GitBranch  = "none"
	// GitMerge   = "0"
	// CiBuild    = "0"
)

func main() {
	app := &cli.App{
		Name:    AppName,
		Usage:   AppUsage,
		Version: AppVersion,
		Action:  actions.Run,
		Flags: []cli.Flag{
			flags.SqlDsnFlag,
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Panic(err)
	}
}
