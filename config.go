package edb

import "fmt"

type (

	// Config database connect config
	Config struct {
		Driver    string
		Host      string
		Port      string
		Database  string
		Username  string
		Password  string
		Charset   string
		Collation string
	}
)

// DriverMysql driver support
const DriverMysql = "mysql"

// DNS return dns string
func (c *Config) DNS() string {
	switch c.Driver {
	case DriverMysql:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset)
	}
	return ""
}
