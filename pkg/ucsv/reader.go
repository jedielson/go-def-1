package ucsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jedielson/bookstore/pkg/database"
	"github.com/jedielson/bookstore/pkg/domain"
)

func ReadFile(filePath string, manager database.DBManager) {

	csvfile, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)
	firstLine := true

	for {

		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if firstLine {
			firstLine = false
			continue
		}

		db := manager.GetDB()

		var author domain.Author
		db.Where("name = ?", record[0]).First(&author)

		if len(author.Name) > 0 {
			continue
		}

		author = domain.Author{
			Name: record[0],
		}

		fmt.Printf("%s\n", author.Name)

		db.Create(&author)
	}
}
