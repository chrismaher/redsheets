package redshift

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	*sql.DB
}

const chunkSize = 1000

func (r *Client) Connect() {
	envVars := []string{"host", "port", "user", "dbname"}
	var connectParams []interface{}
	for _, v := range envVars {
		connectParams = append(connectParams, os.Getenv(v))
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", connectParams...)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	r.DB = db
}

func (r *Client) InsertBulk(schema, table string, values [][]interface{}) {
	insertStr := fmt.Sprintf("INSERT INTO %s.%s VALUES ", schema, table)

	maxlen := 0
	for _, v := range values {
		if len(v) > maxlen {
			maxlen = len(v)
		}
	}

	var vals [][]interface{}
	var params []string
	valueStrings := make([]string, 0, chunkSize)
	valueArgs := make([]interface{}, 0, maxlen*chunkSize)
	var upper, idx int

	for i := 0; i < len(values); i += chunkSize {
		idx = 1
		length := len(values)
		if length < i+chunkSize {
			upper = length
		} else {
			upper = i + chunkSize
		}
		vals = values[i:upper]
		for _, v := range vals {
			params = make([]string, maxlen)
			for j := range v {
				params[j] = "$" + strconv.Itoa(idx)
				valueArgs = append(valueArgs, v[j])
				idx++
			}
			valueStrings = append(valueStrings, "("+strings.Join(params, ", ")+")")
		}
		stmt := insertStr + strings.Join(valueStrings, ", ")

		log.Printf("Inserting into %s.%s %d %d", schema, table, i, idx)
		_, err := r.DB.Exec(stmt, valueArgs...)
		if err != nil {
			log.Println(err)
		}
		valueArgs = valueArgs[:0]
		valueStrings = valueStrings[:0]
	}
}

func (r *Client) Replace(schema, table string, values [][]interface{}) {
	stmt := fmt.Sprintf("TRUNCATE TABLE %s.%s", schema, table)

	log.Printf("Truncating %s.%s", schema, table)
	if _, err := r.DB.Exec(stmt); err != nil {
		panic(err)
	}

	r.InsertBulk(schema, table, values)
}
