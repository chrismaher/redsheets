package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Table struct {
	SheetID   string `json:"sheet_id"`
	SheetName string `json:"sheet_name"`
	Schema    string `json:"schema"`
	Name      string `json:"name"`
}

type Tables map[int]Table

type Data struct {
	Path string
	Map  Tables
}

// Read the JSON file into a slice of Table
func (d *Data) Read() error {
	// test if datastore does not exist yet
	if _, err := os.Stat(d.Path); os.IsNotExist(err) {
		if err := ioutil.WriteFile(d.Path, []byte("{}"), 0644); err != nil {
			return err
		}
		fmt.Printf("%s created\n", d.Path)
	} else {
		b, err := ioutil.ReadFile(d.Path)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(b, &d.Map); err != nil {
			return err
		}
	}

	return nil
}

// Tab-print slice of Table
func (d *Data) List() {
	TabPrinter(d.Map)
}

// Serialize a map[int]Table as to JSON file
func (d *Data) Write() error {
	tableJSON, err := json.Marshal(d.Map)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(d.Path, tableJSON, 0644)
}

// Add a Table to an existing Tables map
func (d *Data) Add(table Table) (int, error) {
	key := len(d.Map)
	d.Map[key] = table
	if err := d.Write(); err != nil {
		return -1, err
	}

	return key, nil
}

// Get a Table by its key from a Tables map
func (d *Data) Get(key int) (Table, error) {
	table, ok := d.Map[key]
	if !ok {
		return Table{}, fmt.Errorf("Table %d does not exist", key)
	}

	return table, nil
}

// Delete a Table by its key from a Tables map
func (d *Data) Delete(key int) error {
	// shift elements at higher indices down one key
	for i := key; i < len(d.Map); i++ {
		d.Map[i] = d.Map[i+1]
	}

	delete(d.Map, len(d.Map)-1)

	if err := d.Write(); err != nil {
		return err
	}

	return nil
}
