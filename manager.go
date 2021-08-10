package edb

import (
	"database/sql"
)

type (

	// Manager magage database connections
	Manager struct {
		connect *connect
	}
)

var manager = &Manager{
	connect: newConnect(),
}

// Exec exec
func (m *Manager) Exec(query string, bindings ...interface{}) (sql.Result, error) {

	return m.connect.Exec(query, bindings...)
}

// Query query
func (m *Manager) Query(query string, bindings ...interface{}) (*sql.Rows, error) {
	return m.connect.Query(query, bindings...)
}

// QueryCollect query and return *Collect
func (m *Manager) QueryCollect(query string, bindings ...interface{}) (*Collect, error) {

	sqlRows, err := m.connect.Query(query, bindings...)
	if err != nil {
		return nil, err
	}
	return &Collect{
		sqlRows: sqlRows,
	}, nil
}
