package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

type Table struct {
	Schema string `json:"schema"`
	Name   string `json:"name"`
}

type Map struct {
	SheetID   string `json:"sheet_id"`
	SheetName string `json:"sheet_name"`
	Table
}

type DataStore struct {
	Path string
	Maps []Map
}

// NewMap takes four string values and returns a Map
func NewMap(id, sheet, schema, name string) Map {
	return Map{id, sheet, Table{schema, name}}
}

// Read the JSON file into a slice of Table
func (d *DataStore) Read() error {
	// test if datastore does not exist yet
	if _, err := os.Stat(d.Path); os.IsNotExist(err) {
		if err := ioutil.WriteFile(d.Path, []byte("[]"), 0644); err != nil {
			return err
		}
		fmt.Printf("%s created\n", d.Path)
	} else {
		b, err := ioutil.ReadFile(d.Path)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(b, &d.Maps); err != nil {
			return err
		}
	}

	return nil
}

// Tab-print slice of Table
func (d *DataStore) Print() {
	str := "%v\t%v\t%v\t%v\t%v\t\n"
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, str, "ID", "Sheet ID", "Sheet Name", "Schema", "Name")
	fmt.Fprintf(tw, str, "--", "--------", "----------", "------", "----")

	var m Map
	for i := 0; i < len(d.Maps); i++ {
		m = d.Maps[i]
		fmt.Fprintf(tw, str, i+1, m.SheetID, m.SheetName, m.Schema, m.Name)
	}

	tw.Flush()
}

// Serialize a []Map to a JSON file
func (d *DataStore) Write() error {
	tableJSON, err := json.Marshal(d.Maps)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(d.Path, tableJSON, 0644)
}

// Add a Table to DataStore.Maps
func (d *DataStore) Add(m Map) (int, error) {
	d.Maps = append(d.Maps, m)
	return len(d.Maps), d.Write()
}

// Get a Table by its key from a Maps map
func (d *DataStore) Get(key int) (Map, error) {
	if key < 1 || key > len(d.Maps) {
		return Map{}, fmt.Errorf("Table %d does not exist", key)
	}

	return d.Maps[key-1], nil
}

// Delete a Table by its key from a Maps map
func (d *DataStore) Delete(key int) error {
	if key < 1 || key > len(d.Maps) {
		return fmt.Errorf("Table %d does not exist", key)
	}

	d.Maps = append(d.Maps[:key-1], d.Maps[key:]...)
	return d.Write()
}
