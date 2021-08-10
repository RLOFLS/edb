package edb

import (
	"testing"
	"time"

	//_ test
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestStmtSelect(t *testing.T) {
	TestBoot(t)

	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"dateTime"`
		UpdatedAt time.Time `type:"dateTime"`
	}

	u := &User{
		Name:      "ttt",
		Age:       111,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m, err := New(u)
	assert.Nil(t, err)
	stmt := m.stmt

	stmt.SetOp(OPSelect)
	err2 := stmt.Build()
	assert.Nil(t, err2)
	assert.Equal(t,
		"SELECT * FROM `user` ;",
		stmt.PrepareSQL(),
	)

	m.reset()
	stmt.SetOp(OPSelect)
	m.builder.Select([]string{"id", "name"})
	stmt.Build()
	assert.Equal(t,
		"SELECT `id`, `name` FROM `user` ;",
		stmt.PrepareSQL(),
	)

	af := func(ecpectSql string, expectbindings []interface{}) {
		assert.Equal(t,
			ecpectSql,
			stmt.PrepareSQL(),
		)
		assert.Equal(t,
			expectbindings,
			stmt.Bindings(),
		)
	}
	m.reset()
	stmt.SetOp(OPSelect)
	m.builder.Select([]string{"name"})
	m.builder.WhereCondition("id", "=", 1)
	m.builder.WhereCondition("age", ">", 20)
	stmt.Build()
	af("SELECT `name` FROM `user` WHERE `id` = ? AND `age` > ? ;", []interface{}{1, 20})

	m.reset()
	stmt.SetOp(OPSelect)
	m.builder.WhereCondition("age", ">", 20)
	m.builder.OrderBy("id")
	m.builder.OrderByDesc("name")
	stmt.Build()
	af("SELECT * FROM `user` WHERE `age` > ? ORDER BY `id` ASC , `name` DESC ;", []interface{}{20})

	m.reset()
	stmt.SetOp(OPSelect)
	m.builder.limit = 1
	m.builder.limitOffset = 0
	stmt.Build()
	af("SELECT * FROM `user` LIMIT 1 ;", []interface{}{})

	m.reset()
	stmt.SetOp(OPSelect)
	m.builder.limit = 1
	stmt.Build()
	af("SELECT * FROM `user` LIMIT 10 OFFSET 0 ;", []interface{}{})

	m.reset()
	stmt.SetOp(OPSelect)
	m.builder.Select([]string{"name"})
	m.builder.WhereCondition("id", "=", 1)
	m.builder.WhereCondition("age", ">", 20)
	m.builder.OrderBy("id")
	m.builder.OrderByDesc("name")
	m.builder.limit = 2
	stmt.Build()
	af("SELECT `name` FROM `user` WHERE `id` = ? AND `age` > ? ORDER BY `id` ASC , `name` DESC LIMIT 10 OFFSET 10 ;", []interface{}{1, 20})

}

func TestStmtUpdate(t *testing.T) {
	TestBoot(t)

	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"dateTime"`
		UpdatedAt time.Time `type:"dateTime"`
	}

	u := &User{
		Name:      "ttt",
		Age:       111,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m, err := New(u)
	assert.Nil(t, err)
	stmt := m.stmt

	stmt.SetOp(OPUpdate)
	err2 := stmt.Build()
	assert.EqualError(t, err2, "edb StmtMysql.Build err: OPUpdate no updated fields")

	m.builder.Update([]string{"name", "age"})
	stmt.Build()
	assert.Equal(t,
		"UPDATE `user` SET `name` = ?,`age` = ? WHERE `id` = ? ;",
		stmt.PrepareSQL(),
	)
	assert.Equal(t,
		[]interface{}{"ttt", 111, 0},
		stmt.Bindings(),
	)

	m.reset()
	stmt.SetOp(OPUpdate)
	m.builder.Update([]string{"age"})
	m.builder.WhereCondition("name", "LIKE", "tom%")
	stmt.Build()
	assert.Equal(t,
		"UPDATE `user` SET `age` = ? WHERE `name` LIKE ? ;",
		stmt.PrepareSQL(),
	)
	assert.Equal(t,
		[]interface{}{111, "tom%"},
		stmt.Bindings(),
	)
}

func TestStmtInsert(t *testing.T) {
	TestBoot(t)

	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"time"`
		UpdatedAt time.Time `type:"dateTime"`
	}

	u := &User{
		Name:      "ttt",
		Age:       111,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m, err := New(u)
	assert.Nil(t, err)
	stmt := m.stmt

	stmt.SetOp(OPInsert)
	stmt.Build()

	//map is an unordered , so the order of the fields cannot be determined
	assert.Equal(t,
		"INSERT INTO `user` (`name`,`age`,`created_at`,`updated_at`) VALUES (?,?,?,?) ;",
		stmt.PrepareSQL(),
	)
	assert.Equal(t,
		[]interface{}{"ttt", 111, "16:22:22", "2021-08-09 16:22:22"},
		stmt.Bindings(),
	)

	//no Auto
	type User2 struct {
		Id        string `type:"pk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"time"`
		UpdatedAt time.Time `type:"dateTime"`
	}
	u2 := &User2{
		Id:        "sss",
		Name:      "ttt",
		Age:       111,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m2, _ := New(u2)
	m2.stmt.SetOp(OPInsert)
	m2.stmt.Build()
	assert.Equal(t,
		"INSERT INTO `user` (`id`,`name`,`age`,`created_at`,`updated_at`) VALUES (?,?,?,?,?) ;",
		m2.stmt.PrepareSQL(),
	)
	assert.Equal(t,
		[]interface{}{"sss", "ttt", 111, "16:22:22", "2021-08-09 16:22:22"},
		m2.stmt.Bindings(),
	)
}

func TestStmtDelete(t *testing.T) {
	TestBoot(t)

	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"dateTime"`
		UpdatedAt time.Time `type:"dateTime"`
	}

	u := &User{
		Name:      "ttt",
		Age:       111,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m, err := New(u)
	assert.Nil(t, err)
	stmt := m.stmt

	stmt.SetOp(OPDelete)
	stmt.Build()
	assert.Equal(t,
		"DELETE FROM `user` WHERE `id` = ? ;",
		stmt.PrepareSQL(),
	)
	assert.Equal(t,
		[]interface{}{0},
		stmt.Bindings(),
	)

	m.reset()
	stmt.SetOp(OPDelete)
	m.Gt("age", 20)
	m.Like("name", "tom%")
	stmt.Build()
	assert.Equal(t,
		"DELETE FROM `user` WHERE `age` > ? AND `name` LIKE ? ;",
		stmt.PrepareSQL(),
	)
	assert.Equal(t,
		[]interface{}{20, "tom%"},
		stmt.Bindings(),
	)
}
