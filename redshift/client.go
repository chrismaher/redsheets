package redshift

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const chunkSize = 1000

func (r *Client) Truncate(schema, table string) error {
	stmt := fmt.Sprintf("TRUNCATE TABLE %s.%s", schema, table)

	log.Printf("Truncating %s.%s", schema, table)
	if _, err := r.Exec(stmt); err != nil {
		return err
	}

	return nil
}

func (r *Client) Insert(vals [][]interface{}, insertStr string, maxlen int, strs *[]string, args *[]interface{}) error {
	idx := 1
	var params []string
	for _, v := range vals {
		params = make([]string, maxlen)
		for j := range v {
			params[j] = "$" + strconv.Itoa(idx)
			*args = append(*args, v[j])
			idx++
		}
		*strs = append(*strs, "("+strings.Join(params, ", ")+")")
		params = params[:0]
	}
	stmt := insertStr + strings.Join(*strs, ", ")

	_, err := r.Exec(stmt, *args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Client) InsertBulk(schema, table string, values [][]interface{}) error {
	insertStr := fmt.Sprintf("INSERT INTO %s.%s VALUES ", schema, table)

	// get the max row width (# of values), so that we can ensure we're always
	// inserting a value for every column
	var maxlen int
	for _, v := range values {
		if len(v) > maxlen {
			maxlen = len(v)
		}
	}

	var valueStrings = make([]string, 0, chunkSize)
	var valueArgs = make([]interface{}, 0, maxlen*chunkSize)

	var length = len(values)
	var upper int

	for i := 0; i < len(values); i += chunkSize {
		fmt.Println(i, len(values), i < len(values))
		if length < i+chunkSize {
			upper = length
		} else {
			upper = i + chunkSize
		}

		vals := values[i:upper]
		log.Printf("Inserting into %s.%s %d", schema, table, i)

		if err := r.Insert(vals, insertStr, maxlen, &valueStrings, &valueArgs); err != nil {
			return err
		}

		// empty the slices before the next iteration
		valueArgs = valueArgs[:0]
		valueStrings = valueStrings[:0]
	}

	return nil
}

func (r *Client) Replace(schema, table string, values [][]interface{}) error {
	if err := r.Truncate(schema, table); err != nil {
		return err
	}

	if err := r.InsertBulk(schema, table, values); err != nil {
		return err
	}

	return nil
}
