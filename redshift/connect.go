package redshift

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type Client struct {
	*sql.DB
}

func (r *Client) Connect() error {
	envVars := []string{"host", "port", "user", "dbname"}
	var connectParams []interface{}
	for _, v := range envVars {
		connectParams = append(connectParams, os.Getenv(v))
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", connectParams...)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	r.DB = db
	return nil
}
