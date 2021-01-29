package actions

import (
	"github.com/jedielson/bookstore/pkg/database"
	"github.com/jedielson/bookstore/pkg/ucsv"
	"github.com/urfave/cli/v2"
)

func Run(c *cli.Context) error {

	manager := database.NewDbManager("bookstore.db")
	err := manager.InitDb()
	if err != nil {
		panic(err)
	}

	database.Migrate(manager.GetDB())

	if err != nil {
		panic(err)
	}

	//file := "../../input.csv"
	file := "./input.csv"
	ucsv.ReadFile(file, manager)

	err = manager.Close()
	if err != nil {
		panic(err)
	}

	return nil
}
