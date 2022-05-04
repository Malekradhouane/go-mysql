package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

//Config db config
type Config struct {
	Host     string
	Port     int
	DB       string
	User     string
	Password string

}


// Client a postgres database client
type Client struct {
	db     *sql.DB
	models []interface{}
}


// NewClient create a new client connecting to a PostgreSQL database
func NewClient(c *Config, models []interface{}) (*Client, error) {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/ring-over-api?charset=utf8")

	if err != nil {
		return nil, errors.Wrap(err, "sql Open")
	}
	defer db.Close()

	return &Client{
		db:     db,
		models: models,
	}, nil
}

// Shutdown close client connection to db
func (c *Client) Shutdown() error {
	return errors.Wrap(c.db.Close(), "Close")
}
