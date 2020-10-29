package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chrismaher/redsheets/homedir"
)

type Table struct {
	SheetID   string `json:"sheet_id"`
	SheetName string `json:"sheet_name"`
	Schema    string `json:"schema"`
	Name      string `json:"name"`
}

type Tables map[int]Table

// Init creates the JSON file in which redsheets will persist Table records
func Init() error {
	db, err := homedir.FullPath(".redsheets.json")
	if err != nil {
		return err
	}

	if _, err := os.Stat(db); os.IsNotExist(err) {
		if err := ioutil.WriteFile(db, []byte("{}"), 0644); err != nil {
			return err
		}
		fmt.Printf("%s created\n", db)
		return nil
	}

	return fmt.Errorf("%s already exists", db)
}

func parseTables(b []byte) (Tables, error) {
	var data Tables
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// Read the JSON file into a slice of Table
func Read() (Tables, error) {
	db, err := homedir.FullPath(".redsheets.json")
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

	TabPrinter(tables)

	return nil
}

func writeJSON(data Tables) error {
	db, err := homedir.FullPath(".redsheets.json")
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
		return err
	}

	for _, v := range data {
		if table.SheetID == v.SheetID {
			return fmt.Errorf("Sheet ID %s is already in db", table.SheetID)
		}
	}

	data[len(data)] = table
	if err = writeJSON(data); err != nil {
		return err
	}

	return nil
}

func Get(key int) (Table, error) {
	tables, err := Read()
	if err != nil {
		return Table{}, err
	}

	table, ok := tables[key]
	if !ok {
		return Table{}, fmt.Errorf("Key %d does not exist", key)
	}

	return table, nil
}

func Delete(index int) error {
	data, err := Read()
	if err != nil {
		return err
	}

	delete(data, index)

	// shift elements at higher indices down one index
	for i := index + 1; i < len(data)+1; i++ {
		data[i-1] = data[i]
	}

	delete(data, len(data)-1)

	if err = writeJSON(data); err != nil {
		return err
	}

	return nil
}
