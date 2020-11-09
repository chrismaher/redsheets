package redshift

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Connection struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Name string `mapstructure:"dbname"`
}

type Client struct {
	*Connection
	*sql.DB
}

func (c *Connection) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Name)
}

func (r *Client) Connect() error {
	db, err := sql.Open("postgres", r.String())
	if err != nil {
		return err
	}

	r.DB = db
	return nil
}
