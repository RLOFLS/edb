package edb

import "database/sql"

type (

	// operateType oprate types
	// select delete update insert
	operateType int

	// Stmt stmt
	Stmt interface {
		// Build build sql, and bind parameters
		Build() error
		// PrepareSQL get prepare sql
		PrepareSQL() string
		// Bindings get prepare bind parameters
		Bindings() []interface{}
		// SetBuilder set Builder
		SetBuilder(*Builder)
		// SetOp
		SetOp(operateType)
		reset()
	}

	// Query query
	Query interface {
		// Select
		Select([]string) *Model
		// Eq("name", "tom") => name='tom'
		Eq(field string, value interface{}) *Model
		// Neq("name", "tom") => name !='tom'
		Neq(field string, value interface{}) *Model
		// Lt("age", 1) => `age` < 1
		Lt(field string, value interface{}) *Model
		// Lte("age", 1) => `age` <= 1
		Lte(field string, value interface{}) *Model
		// Gt("age", 1) => `age` > 1
		Gt(field string, value interface{}) *Model
		// Gte("age", 1) => `age` >= 1
		Gte(field string, value interface{}) *Model
		// Like("name", "%sss") => `name` LIKE '%sss'
		Like(field string, value interface{}) *Model
		// Order ASC sort
		OrderBy(string) *Model
		// OrderByDesc DESC sort
		OrderByDesc(string) *Model
		Get() (*Collect, error)
		First() (interface{}, error)
		Paginate(page int64, pageSize int64) (*Collect, error)
		Delete() (rowAffected int64, err error)
		Insert() (id int64, err error)
		Update([]string) (rowAffected int64, err error)
		// todo
		// Transaction(conn string, closure func()) error
		Query(string, ...interface{}) (*sql.Rows, error)
		Exec(string, ...interface{}) (sql.Result, error)
		// todo
		// ToSQL() string
	}
)
