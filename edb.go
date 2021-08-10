/*
Package edb simple and convenient operation of database data table, structure mapping

Example:

  package main

  import (
  	"github.com/RLOFLS/edb"
  	"fmt"
  	"log"
  	"time"

  	_ "github.com/go-sql-driver/mysql"
  )

  func main() {

    //start
    edb.AddConfig("default", &edb.Config{Driver: "mysql", Host: "127.0.0.1", Port: "3306", Username: "root", Password: "12345678", Database: "test", Charset: "utf8"})
    edb.Boot("default")

    //define structure corresponds to the data table
    //naming rules:
    //database : user_config   => struct : UserConfig
    type User struct {
    	Id        int `type:"autoPk"`
    	Name      string
    	Age       int
    	CreatedAt time.Time `type:"dateTime"`
    	UpdatedAt time.Time `type:"dateTime"`
    }

    lf := func(err error) {
    	if err != nil {
    		log.Fatal(err)
    	}
    }

    //Insert
    fmt.Println("---Insert----")
    t, _ := time.ParseInLocation(edb.FTimeDateTime, "2021-01-01 01:01:01", time.Local)
    m, err := edb.New(&User{
    	Name:      "Test",
    	Age:       222,
    	CreatedAt: t,
    	UpdatedAt: t,
    })
    lf(err)
    id, err := m.Insert()
    lf(err)
    fmt.Println("last insert id:", id)

    //Select support Eq, Lt, Lte, Gt, Gte, Like, OrderBy, OrderByDesc chain opreation
    // eg： m.Eq("name", "a").Gt("age", 22)
    fmt.Println("---Select----")
    m2, err := edb.New(&User{})
    lf(err)
    i2, err := m2.Eq("name", "test").First()
    lf(err)
    user2 := i2.(*User)
    fmt.Printf("Select user： userName: %s, age: %d, createTime: %s\n", user2.Name, user2.Age, user2.CreatedAt.Format(edb.FTimeDateTime))

    //Get all
    fmt.Println("---Get----")
    m3, err := edb.New(&User{})
    lf(err)
    collect3, err := m3.Eq("name", "test").Get()
    lf(err)
    for collect3.Next() {
    	u := collect3.Item()
    	user3 := u.(*User)
    	fmt.Printf("Get User: userName: %s, age: %d, createTime: %s\n", user3.Name, user3.Age, user3.CreatedAt.Format(edb.FTimeDateTime))
    }

    //Paginate
    fmt.Println("---Paginate----")
    m4, err := edb.New(&User{})
    lf(err)
    collect4, err := m4.Like("name", "%t%").OrderByDesc("created_at").Paginate(1, 10)
    lf(err)
    fmt.Printf("Paginate User: total count: %d\n", collect4.Total())
    for collect4.Next() {
    	u := collect4.Item()
    	user4 := u.(*User)
    	fmt.Printf("Paginate User: userName: %s, age: %d, createTime: %s\n", user4.Name, user4.Age, user4.CreatedAt.Format(edb.FTimeDateTime))
    }

    //Update
    fmt.Println("---Update----")
    m5, err := edb.New(&User{
    	Id:        1,
    	Name:      "ttttt",
    	Age:       111,
    	UpdatedAt: time.Now(),
    })
    lf(err)
    //default pk as where condition => where id=1
    rowAffected, err := m5.Update([]string{"name", "age", "created_at"})
    lf(err)
    fmt.Printf("Update User: rowAffected : %d\n", rowAffected)

    //Delete
    fmt.Println("---Delete----")
    m6, err := edb.New(&User{})
    lf(err)
    rowAffected2, err := m6.Eq("name", "ttttt").Delete()
    lf(err)
    fmt.Printf("Delete User: rowAffected : %d\n", rowAffected2)

    //customize structure mapping, just query
    fmt.Println("---customize structure mapping----")
    type UserCount struct {
    	UserName string
    	Total    int
    }
    m7, err := edb.New(&UserCount{})
    lf(err)
    prepareSQL := "SELECT `name` AS user_name, count(*) AS total FROM `user` GROUP BY `name`"
    collect7, err := m7.QueryCollect(prepareSQL)
    lf(err)
    for collect7.Next() {
    	i := collect7.Item()
    	userCount := i.(*UserCount)
    	fmt.Printf("userName: %s, count: %d\n", userCount.UserName, userCount.Total)
    }

    //and more usage see package test file

  }
*/
package edb

import "fmt"

const (
	// OPSelect select
	OPSelect = iota
	// OPDelete delete
	OPDelete
	// OPUpdate update
	OPUpdate
	// OPInsert insert
	OPInsert
)

// structrue tag
// time formatting
const (
	TagAutoPK   = "autoPk"
	TagPK       = "pk"
	TagDate     = "date"
	TagDateTime = "dateTime"
	TagTime     = "time"

	FTimeTime     = "15:04:05"
	FTimeDate     = "2006-01-02"
	FTimeDateTime = "2006-01-02 15:04:05"
)

// AddConfig add database connection configuration
func AddConfig(connectName string, config *Config) bool {
	return manager.connect.AddConfig(connectName, config)
}

// Boot start up
func Boot(connectName string) {
	if err := manager.connect.Connect(connectName); err != nil {
		panic(err.Error())
	}
}

func newStmt() (stmt Stmt, err error) {
	switch manager.connect.driver {
	case DriverMysql:
		stmt = &StmtMysql{
			bindings: make([]interface{}, 0),
		}
	default:
		err = fmt.Errorf("edb Model.New err: unsupported %s driver syntax", manager.connect.driver)
	}
	return
}
