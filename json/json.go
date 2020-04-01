package json

import (
	"encoding/json"
	"fmt"
	"github.com/chrismaher/redsheets/path"
	"io/ioutil"
	"log"
	"os"
)

type Table struct {
	SheetID   string `json:"sheet_id"`
	SheetName string `json:"sheet_name"`
	Schema    string `json:"schema"`
	Name      string `json:"name"`
}

// Init creates the JSON file in which redsheets will persist Table records
func Init() error {
	db, err := path.FullPath(".redsheets.json")
	if err != nil {
		return err
	}

	if _, err := os.Stat(db); os.IsNotExist(err) {
		if err := ioutil.WriteFile(db, []byte("[]"), 0644); err != nil {
			return err
		}
		log.Printf("%s created", db)
		return nil
	}

	return fmt.Errorf("%s already exists", db)
}

func parseTables(b []byte) ([]Table, error) {
	var data []Table
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// Read the JSON file into a slice of Table
func Read() ([]Table, error) {
	db, err := path.FullPath(".redsheets.json")
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(db)
	if err != nil {
		return nil, err
	}

	return parseTables(b)
}

func List() error {
	tables, err := Read()
	if err != nil {
		return err
	}

	for _, table := range tables {
		log.Printf("%+v", table)
	}

	return nil
}

func writeJSON(data []Table) error {
	db, err := path.FullPath(".redsheets.json")
	if err != nil {
		return err
	}

	tableJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(db, tableJSON, 0644)
}

func Add(table Table) error {
	data, err := Read()
	if err != nil {
		return nil
	}

	for _, t := range data {
		if table.SheetID == t.SheetID {
			return fmt.Errorf("Sheet ID %s is already in db", table.SheetID)
		}
	}

	if err = writeJSON(append(data, table)); err != nil {
		return err
	}

	log.Println("Added")
	return nil
}
