package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

type Table struct {
	SheetID   string `json:"sheet_id"`
	SheetName string `json:"sheet_name"`
	Schema    string `json:"schema"`
	Name      string `json:"name"`
}

type TableMap map[int]Table

type DataStore struct {
	Path   string
	Tables TableMap
}

// Read the JSON file into a slice of Table
func (d *DataStore) Read() error {
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

		if err := json.Unmarshal(b, &d.Tables); err != nil {
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

	var v Table
	for i := 0; i < len(d.Tables); i++ {
		v = d.Tables[i]
		fmt.Fprintf(tw, str, i, v.SheetID, v.SheetName, v.Schema, v.Name)
	}

	tw.Flush()
}

// Serialize a TableMap to a JSON file
func (d *DataStore) Write() error {
	tableJSON, err := json.Marshal(d.Tables)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(d.Path, tableJSON, 0644)
}

// Add a Table to an existing Tables map
func (d *DataStore) Add(table Table) (int, error) {
	key := len(d.Tables)
	d.Tables[key] = table
	if err := d.Write(); err != nil {
		return -1, err
	}

	return key, nil
}

// Get a Table by its key from a Tables map
func (d *DataStore) Get(key int) (Table, error) {
	table, ok := d.Tables[key]
	if !ok {
		return Table{}, fmt.Errorf("Table %d does not exist", key)
	}

	return table, nil
}

// Delete a Table by its key from a Tables map
func (d *DataStore) Delete(key int) error {
	// shift elements at higher indices down one key
	for i := key; i < len(d.Tables); i++ {
		d.Tables[i] = d.Tables[i+1]
	}

	delete(d.Tables, len(d.Tables)-1)

	if err := d.Write(); err != nil {
		return err
	}

	return nil
}
