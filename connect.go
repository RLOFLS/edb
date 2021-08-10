package edb

import (
	"database/sql"
	"fmt"
)

type (

	// connect connections
	connect struct {
		db     *sql.DB
		config map[string]*Config
		driver string
	}
)

// newConnect new
func newConnect() *connect {
	return &connect{
		config: make(map[string]*Config, 10),
	}
}

// AddConfig add connection configuration
func (conn *connect) AddConfig(connectName string, config *Config) bool {
	if config.DNS() == "" {
		return false
	}
	conn.config[connectName] = config
	return true
}

// Connect connect to the database
func (conn *connect) Connect(connectName string) error {

	var (
		c  *Config
		ok bool
	)

	if c, ok = conn.config[connectName]; !ok {
		return fmt.Errorf("edb: connect.Connect err: %s : database connection configuration dose not exist", connectName)
	}

	conn.driver = c.Driver

	db, err := sql.Open(c.Driver, c.DNS())
	if err != nil {
		return fmt.Errorf("edb: connect.Connect err: %s", err.Error())
	}

	err2 := db.Ping()
	if err2 != nil {
		return fmt.Errorf("edb: connect.Connect err: %s", err2.Error())
	}
	conn.db = db

	return nil
}

// Db sql.DB
func (conn *connect) Db() *sql.DB {
	return conn.db
}

// Exec DB.Exec
func (conn *connect) Exec(query string, args ...interface{}) (sql.Result, error) {
	return conn.db.Exec(query, args...)
}

// Query query
func (conn *connect) Query(query string, args ...interface{}) (*sql.Rows, error) {

	stmt, err := conn.db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	return stmt.Query(args...)
}
