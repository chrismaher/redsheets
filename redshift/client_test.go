package redshift

import "testing"

func TestConnect(t *testing.T) {
	db := Client{}
	db.Connect()
	defer db.DB.Close()

	err := db.DB.Ping()
	if err != nil {
		t.Errorf("Could not connect to specified database")
	}
}
