package edb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	//_ test
	_ "github.com/go-sql-driver/mysql"
)

func TestAddConfig(t *testing.T) {
	c := newConnect()
	c.AddConfig("test", &Config{})
}

func TestConnect(t *testing.T) {
	c := newConnect()
	err := c.Connect("test")
	assert.Equal(t, "edb: connect.Connect err: test : database connection configuration dose not exist", err.Error())

	c.AddConfig("test", &Config{Driver: "mysql", Host: "127.0.0.1", Port: "3306", Username: "root", Password: "12345678", Database: "test", Charset: "utf8"})
	err = c.Connect("test")
	assert.NoError(t, err)
}

func TestExec(t *testing.T) {
	c := newConnect()
	c.AddConfig("test", &Config{Driver: "mysql", Host: "127.0.0.1", Port: "3306", Username: "root", Password: "12345678", Database: "test", Charset: "utf8"})
	err := c.Connect("test")
	assert.NoError(t, err)

	tests := []struct {
		query        string
		rowsAffected int64
	}{
		{
			"DROP TABLE IF EXISTS `user`;",
			0,
		},
		{
			"CREATE TABLE `test`.`user`(`id` int(0) AUTO_INCREMENT NOT NULL,`name` varchar(50) NULL DEFAULT '', `age` int unsigned DEFAULT '0', `created_at` datetime DEFAULT NULL, `updated_at` datetime DEFAULT NULL,PRIMARY KEY (`id`));",
			0,
		},
		{
			"INSERT INTO user (`name`) values('tom')",
			1,
		},
	}

	for _, test := range tests {
		res, err := c.Exec(test.query)
		assert.NoError(t, err)
		ra, err := res.RowsAffected()
		assert.Equal(t, ra, test.rowsAffected)
	}
}

func TestQuery(t *testing.T) {
	c := newConnect()
	c.AddConfig("test", &Config{Driver: "mysql", Host: "127.0.0.1", Port: "3306", Username: "root", Password: "12345678", Database: "test", Charset: "utf8"})
	err := c.Connect("test")
	assert.NoError(t, err)

	p := "SELECT id, name FROM `user`"
	sqlRows, err := c.Query(p)
	defer sqlRows.Close()
	assert.NoError(t, err)

	if sqlRows.Next() {
		cls, err := sqlRows.Columns()
		if err != nil {
			t.Error(err)
		}
		t.Log(cls)

		c := new(string)
		d := new(string)
		sqlRows.Scan(c, d)

		assert.Equal(t, "tom", *d)

	}
}
