package edb

import (
	"database/sql"
	"errors"
	"reflect"
	"time"
)

type (

	// Collect collect
	Collect struct {
		currentEntity interface{}
		sqlRows       *sql.Rows
		paginateTotal int64
		originModel   *Model
		err           error
	}
)

// Item get the value of the current iteration mapping
func (c *Collect) Item() interface{} {
	if c.currentEntity != nil {
		return c.currentEntity
	}
	if c.Next() {
		return c.currentEntity
	}
	return nil
}

// Next iterative data, if return true , call Item() get value
// calling again will iterate the next item of data
func (c *Collect) Next() bool {
	fn := func() bool {
		if err := c.setCurrentEntity(); err != nil {
			c.err = err
			return false
		}
		return true
	}
	if c.sqlRows.Next() {
		return fn()
	}
	//no data
	c.sqlRows.Close()
	return false
}

// Total get the total number of queries, just Model.Paginate
func (c *Collect) Total() int64 {
	return c.paginateTotal
}

// Err get err
func (c *Collect) Err() error {
	return c.err
}

// setCurrentModel set currentEntity
func (c *Collect) setCurrentEntity() error {

	if c.originModel == nil {
		return errors.New("edb Collect.setCurrentEntity err: attr originModel is not set")
	}

	//new from original structure
	rValue := reflect.New(reflect.TypeOf(c.originModel.entity).Elem())

	c.currentEntity = rValue.Interface()

	columns, err := c.sqlRows.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	c.initScanValues(columns, values)
	if err := c.sqlRows.Scan(values...); err != nil {
		return err
	}

	//assign
	for idx, cName := range columns {
		if field, ok := c.originModel.entityFields[cName]; ok {
			fValue := rValue.Elem().FieldByName(field.sName)
			if fValue.CanSet() {
				switch field.fType {
				case "string":
					fValue.SetString(*values[idx].(*string))
				case "int", "int8", "int16", "int32", "int64":
					fValue.SetInt(*values[idx].(*int64))
				case "uint", "uint8", "uint16", "uint32", "uint64":
					fValue.SetUint(*values[idx].(*uint64))
				case "bool":
					fValue.SetBool(*values[idx].(*bool))
				case "float32", "float64":
					fValue.SetFloat(*values[idx].(*float64))
				case "Time", "time.Time":
					switch field.fTagType {
					case TagTime:
						if t, err := time.ParseInLocation(FTimeTime, *values[idx].(*string), time.Local); err == nil {
							fValue.Set(reflect.ValueOf(t))
						}
					case TagDate:
						if t, err := time.ParseInLocation(FTimeDate, *values[idx].(*string), time.Local); err == nil {
							fValue.Set(reflect.ValueOf(t))
						}
					case TagDateTime:
						if t, err := time.ParseInLocation(FTimeDateTime, *values[idx].(*string), time.Local); err == nil {
							fValue.Set(reflect.ValueOf(t))
						}
					}
				default:
				}
			}

		}

	}
	return nil
}

func (c *Collect) initScanValues(columns []string, values []interface{}) {
	for idx, cName := range columns {
		field, ok := c.originModel.entityFields[cName]
		if !ok {
			values[idx] = new(interface{})
		} else {
			switch field.fType {
			case "string":
				values[idx] = new(string)
			case "int", "int8", "int16", "int32", "int64":
				values[idx] = new(int64)
			case "uint", "uint8", "uint16", "uint32", "uint64":
				values[idx] = new(uint64)
			case "bool":
				values[idx] = new(bool)
			case "float32", "float64":
				values[idx] = new(float64)
			case "Time", "time.Time":
				values[idx] = new(string)
			default:
				values[idx] = new(interface{})
			}
		}
	}
}
