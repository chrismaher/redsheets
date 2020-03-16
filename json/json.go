package json

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

type Table struct {
	SheetID   string `json:"sheet_id"`
	SheetName string `json:"sheet_name"`
	Schema    string `json:"schema"`
	Name      string `json:"name"`
}

func FullPath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}

	home := usr.HomeDir
	return path.Join(home, filename)
}

var db = FullPath(".redsheets.json")

func Init() {
	if _, err := os.Stat(db); os.IsNotExist(err) {
		err := ioutil.WriteFile(db, []byte("[]"), 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s created", db)
		return
	}
	log.Printf("%s already exists", db)
}

func Read() []Table {
	b, err := ioutil.ReadFile(db)
	if err != nil {
		panic(err)
	}
	var data []Table
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func List() {
	tables := Read()
	for _, table := range tables {
		log.Printf("%+v", table)
	}
}

func Add(table Table) {
	data := Read()
	for _, t := range data {
		if table.SheetID == t.SheetID {
			log.Printf("Sheet ID %s is already in db", table.SheetID)
			return
		}
	}
	updated, _ := json.Marshal(append(data, table))
	err := ioutil.WriteFile(db, updated, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Added")
}
