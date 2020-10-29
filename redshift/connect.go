package redshift

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Client struct {
	*sql.DB
}

func (r *Client) Connect() error {
	var connectParams []interface{}
	for _, v := range []string{"host", "port", "user", "dbname"} {
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
