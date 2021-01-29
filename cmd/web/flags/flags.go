package flags

import "github.com/urfave/cli/v2"

var (
	SqlDsnFlag = &cli.StringFlag{
		Name:     "sql-dsn",
		Usage:    "dsn to use for connecting database",
		Value:    "lol",
		EnvVars:  []string{"BOOKSTORE_SQL_DSN"},
		Required: false,
	}
)
