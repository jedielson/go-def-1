package actions

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jedielson/bookstore/pkg/api"
	"github.com/jedielson/bookstore/pkg/database"
	"github.com/urfave/cli/v2"
)

func Run(c *cli.Context) error {

	fmt.Printf("Starting api...\n")
	manager := database.NewDbManager("bookstore.db")
	authorsRepository := database.NewAuthorsRepository(manager)

	err := manager.InitDb()
	if err != nil {
		panic(err)
	}

	database.Migrate(manager.GetDB())

	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	api.NewAuthorsApi(r, authorsRepository)

	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	log.Fatal(http.ListenAndServe(":8081", r))

	err = manager.Close()
	if err != nil {
		panic(err)
	}

	return nil
}
