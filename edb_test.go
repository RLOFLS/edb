package edb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	//_ test
	_ "github.com/go-sql-driver/mysql"
)

func TestEdbAddConfig(t *testing.T) {
	res := AddConfig("default", &Config{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Port:     "3306",
		Username: "root",
		Password: "123",
		Database: "test",
		Charset:  "utf8",
	})
	_, ok := manager.connect.config["default"]
	assert.True(t, ok)
	assert.True(t, res)

	res2 := AddConfig("test", &Config{Host: "127.0.0.1", Port: "3306", Username: "root", Password: "123", Database: "test", Charset: "utf8"})
	_, ok2 := manager.connect.config["test"]
	assert.False(t, ok2)
	assert.False(t, res2)
}

func TestBoot(t *testing.T) {
	res := AddConfig("default", &Config{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Port:     "3306",
		Username: "root",
		Password: "123",
		Database: "test",
		Charset:  "utf8",
	})
	assert.True(t, res)
	assert.PanicsWithValue(t, "edb: connect.Connect err: Error 1045: Access denied for user 'root'@'localhost' (using password: YES)", func() {
		Boot("default")
	})

	res2 := AddConfig("default", &Config{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Port:     "3306",
		Username: "root",
		Password: "12345678",
		Database: "test",
		Charset:  "utf8",
	})
	assert.True(t, res2)
	assert.NotPanics(t, func() {
		Boot("default")
	})
}
