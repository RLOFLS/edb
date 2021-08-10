package edb

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

var edbManager *Manager

type (

	// Model model
	Model struct {
		builder      *Builder
		entity       interface{}
		name         string
		tableName    string
		entityFields map[string]Field
		pkField      string
		isAuto       bool
		lastErr      error
		stmt         Stmt
		//todo
		isChange bool
	}

	// Field field
	Field struct {
		isPk     bool
		isAuto   bool
		name     string
		value    interface{}
		fType    string
		fTagType string
		sName    string
	}

	// SupportTypes support driver syntax,
	// unsupported types will not be processed
	SupportTypes map[string]bool
)

var supportTypes = SupportTypes{
	"string":    true,
	"int":       true,
	"int8":      true,
	"int16":     true,
	"int32":     true,
	"int64":     true,
	"uint":      true,
	"uint8":     true,
	"uint16":    true,
	"uint32":    true,
	"uint64":    true,
	"bool":      true,
	"float32":   true,
	"float64":   true,
	"Time":      true,
	"time.Time": true,
}

var _ Query = &Model{}

// New new model, and init someting
func New(entity interface{}) (m *Model, err error) {
	m = &Model{
		builder:      NewBuilder(),
		entity:       entity,
		entityFields: make(map[string]Field, 50),
	}
	m.builder.model = m
	m.stmt, err = newStmt()
	if err != nil {
		return
	}
	m.stmt.SetBuilder(m.builder)

	if err = m.checkEntity(); err != nil {
		return
	}

	if err = m.setTableAttributes(); err != nil {
		return
	}

	return
}

// Select select fields
func (m *Model) Select(s []string) *Model {
	if err := m.builder.Select(s); err != nil {
		m.lastErr = err
	}
	return m
}

// Eq Eq("name", "tom") => name='tom'
func (m *Model) Eq(field string, value interface{}) *Model {
	if err := m.builder.WhereCondition(field, "=", value); err != nil {
		m.lastErr = err
	}
	return m
}

// Neq Neq("name", "tom") => name !='tom'
func (m *Model) Neq(field string, value interface{}) *Model {
	if err := m.builder.WhereCondition(field, "!=", value); err != nil {
		m.lastErr = err
	}
	return m
}

// Lt Lt("age", 1) => `age` < 1
func (m *Model) Lt(field string, value interface{}) *Model {
	if err := m.builder.WhereCondition(field, "<", value); err != nil {
		m.lastErr = err
	}
	return m
}

// Lte Lte("age", 1) => `age` <= 1
func (m *Model) Lte(field string, value interface{}) *Model {
	if err := m.builder.WhereCondition(field, "<=", value); err != nil {
		m.lastErr = err
	}
	return m
}

// Gt Gt("age", 1) => `age` > 1
func (m *Model) Gt(field string, value interface{}) *Model {
	if err := m.builder.WhereCondition(field, ">", value); err != nil {
		m.lastErr = err
	}
	return m
}

// Gte Gte("age", 1) => `age` >= 1
func (m *Model) Gte(field string, value interface{}) *Model {
	if err := m.builder.WhereCondition(field, ">=", value); err != nil {
		m.lastErr = err
	}
	return m
}

// Like Like("name", "%sss") => `name` LIKE '%sss'
func (m *Model) Like(field string, value interface{}) *Model {
	if err := m.builder.WhereCondition(field, "LIKE", value); err != nil {
		m.lastErr = err
	}
	return m
}

// OrderBy ASC sort
func (m *Model) OrderBy(field string) *Model {
	m.builder.OrderBy(field)
	return m
}

// OrderByDesc DESC sort
func (m *Model) OrderByDesc(field string) *Model {
	m.builder.OrderByDesc(field)
	return m
}

// First get the first, if there is no where condition, the pk will be used as the query condition
func (m *Model) First() (interface{}, error) {
	defer m.reset()

	m.builder.limit = 1
	m.builder.limitOffset = 0
	collect, err := m.querySQL()
	if err != nil {
		return nil, err
	}
	collect.originModel = m
	return collect.Item(), nil
}

// Get get all, return *Collect if there is no where condition, the pk will be used as the query condition，
// take out througth for loop
//
// Example usage:
//
// (
// 	c, err := m.Get()
// 	if err != nil {
//   //...
//  }
//  for c.Next() {
//	 i := c.Item()
//	 //...
//  }
// )
func (m *Model) Get() (collect *Collect, err error) {
	defer m.reset()

	return m.returnCollect()
}

// Paginate paginate query, return *Collect, usage as Get()
func (m *Model) Paginate(page int64, pageSize int64) (collect *Collect, err error) {
	defer m.reset()

	m.builder.limit = page
	m.builder.limitOffset = pageSize
	collect, err = m.returnCollect()

	prepareSQL := m.stmt.PrepareSQL()
	if i := strings.Index(prepareSQL, "FROM"); i >= 0 {
		totalSQL := "SELECT count(*) as paginate " + string([]rune(prepareSQL)[i:])
		sqlRows, err := m.Query(totalSQL, m.stmt.Bindings()...)
		if err != nil {
			return nil, err
		}
		if sqlRows.Next() {
			c := 0
			sqlRows.Scan(&c)
			collect.paginateTotal = int64(c)
		}
		sqlRows.Close()
	}
	return
}

// Delete delete, if there is no where condition, the pk will be used as the query condition
//
// Example usage:
// (
// 	m, err := New(&User{
// 		Id: 1,
// 	})
// 	//default  pk    => where `id` = 1
// 	rowAffected, err := m.Delete()
// )
func (m *Model) Delete() (rowAffected int64, err error) {
	defer m.reset()

	m.stmt.SetOp(OPDelete)
	return m.returnRowAffected()
}

// Update update according to rhe passed field, if there is no where condition, the pk will be used as the query condition
//
// Example usage:
//
// (
// 	m, err := New(&User{
// 		Id:   1,
// 		Name: "ttt",
// 		Age:  333,
// 	})
// 	rowAffected, err := m.Update([]string{"name", "age"})
// )
func (m *Model) Update(updateFields []string) (rowAffected int64, err error) {
	defer m.reset()

	m.stmt.SetOp(OPUpdate)
	m.builder.Update(updateFields)
	return m.returnRowAffected()
}

// Insert insert entity
func (m *Model) Insert() (id int64, err error) {
	defer m.reset()

	m.stmt.SetOp(OPInsert)
	return m.returnLastInsertId()
}

// Transaction todo
// conn string connect key
// closure exec closure
// func (m *model) Transaction(conn string, closure func()) error {
// 	return nil
// }

// Query query and return *sql.Rows
func (m *Model) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return manager.Query(query, args...)
}

// Exec exec and return sql.Resqult
func (m *Model) Exec(query string, args ...interface{}) (sql.Result, error) {
	return manager.Exec(query, args...)
}

// QueryCollect query and return *Collect
func (m *Model) QueryCollect(query string, args ...interface{}) (collect *Collect, err error) {
	collect, err = manager.QueryCollect(query, args...)
	if err != nil {
		return nil, err
	}
	collect.originModel = m
	return
}

// ToSQL todo return sql string
// func (m *Model) ToSQL() string {
// 	return ""
// }

func (m *Model) returnRowAffected() (int64, error) {
	sqlResult, err := m.execSQL()
	if err != nil {
		return 0, err
	}

	return sqlResult.RowsAffected()
}

func (m *Model) returnLastInsertId() (int64, error) {
	sqlResult, err := m.execSQL()
	if err != nil {
		return 0, err
	}

	return sqlResult.LastInsertId()
}

func (m *Model) returnCollect() (collect *Collect, err error) {
	collect, err = m.querySQL()
	if err != nil {
		return nil, err
	}
	collect.originModel = m
	return
}

func (m *Model) checkFinalErrWithRun() error {
	if m.lastErr != nil {
		return m.lastErr
	}

	if err := m.stmt.Build(); err != nil {
		return err
	}
	return nil
}

func (m *Model) execSQL() (sql.Result, error) {

	if err := m.checkFinalErrWithRun(); err != nil {
		return nil, err
	}
	return manager.Exec(m.stmt.PrepareSQL(), m.stmt.Bindings()...)
}

func (m *Model) querySQL() (*Collect, error) {

	if err := m.checkFinalErrWithRun(); err != nil {
		return nil, err
	}
	return manager.QueryCollect(m.stmt.PrepareSQL(), m.stmt.Bindings()...)
}

////struct tag
//`type:"auto_pk"`
//`type:"pk"`
//`type:"date"`
//`type:"dateTime"`
func (m *Model) setTableAttributes() error {
	rv := reflect.ValueOf(m.entity).Elem()
	rvt := rv.Type()
	tableName := rvt.Name()
	fieldNums := rv.NumField()
	for i := 0; i < fieldNums; i++ {

		fType := rv.Field(i).Type().String()
		if _, ok := supportTypes[fType]; !ok {
			return fmt.Errorf("edb Model.setTableAttributes err: unsupported field type: %s", fType)
		}

		fName := rvt.Field(i).Name
		fTagType := rvt.Field(i).Tag.Get("type")
		fDBName := camelToUnerline(fName)

		if !rv.Field(i).CanInterface() {
			return fmt.Errorf("edb Model.setTableAttributes err: field: %s unexported", fName)
		}

		f := Field{
			fType:    fType,
			fTagType: fTagType,
			name:     fDBName,
			sName:    fName,
			value:    rv.Field(i).Interface(),
		}

		autoPK := fTagType == TagAutoPK
		pk := fTagType == TagPK
		if autoPK || pk {
			if m.pkField != "" {
				return fmt.Errorf("edb Model.setTableAttributes err: has been set pk: %s, can no longer set the filed `%s` as pk", m.pkField, fDBName)
			}
			m.pkField = fDBName
			f.isPk = true

			if autoPK {
				m.isAuto = true
				f.isAuto = true
			}
		}
		m.entityFields[fDBName] = f

	}
	m.tableName = camelToUnerline(tableName)
	return nil
}

func (m *Model) checkEntity() error {
	rt := reflect.TypeOf(m.entity)
	if rt.Kind().String() != "ptr" {
		return fmt.Errorf("edb Model.checkEntity err: the parameter \"entity\" must be a struct pointer")
	}
	if rt.Elem().Kind().String() != "struct" {
		return fmt.Errorf("edb Model.checkEntity err: the parameter \"entity\" must be a struct")
	}
	return nil
}

func (m *Model) reset() {
	m.builder.reset()
	m.stmt.reset()
}

func camelToUnerline(s string) string {
	buffer := &strings.Builder{}
	for i, v := range s {
		if unicode.IsUpper(v) {
			if i != 0 {
				buffer.WriteByte('_')
			}
			buffer.WriteRune(unicode.ToLower(v))
		} else {
			buffer.WriteRune(v)
		}
	}
	return buffer.String()
}

func underlineToCamel(s string) string {
	res := make([]byte, 50)
	ss := strings.Split(s, "_")
	for _, v := range ss {
		for i2, v2 := range v {
			if i2 == 0 {
				res = append(res, byte(unicode.ToLower(v2)))
			} else {
				res = append(res, byte(v2))
			}
		}
	}
	return string(res)
}
