package edb

import (
	"fmt"
	"strings"
	"time"
)

type (

	// StmtMysql nysql stmt
	StmtMysql struct {
		builder    *Builder
		prepareSQL string
		bindings   []interface{}
		op         operateType
	}
)

var _ Stmt = &StmtMysql{}

// SetBuilder SetBuilder
func (sm *StmtMysql) SetBuilder(b *Builder) {
	sm.builder = b
}

// SetOp SetOp
func (sm *StmtMysql) SetOp(op operateType) {
	sm.op = op
}

// Build build dql, and bind parameters
func (sm *StmtMysql) Build() error {
	sqlBuffer := new(strings.Builder)
	switch sm.op {
	case OPSelect:
		sqlBuffer.WriteString("SELECT ")
		//builder.fields
		fieldsLen := len(sm.builder.fields)
		if fieldsLen == 0 {
			sqlBuffer.WriteString("* ")
		} else {
			for i, f := range sm.builder.fields {
				if fieldsLen == i+1 {
					sqlBuffer.WriteString("`" + f + "` ")
				} else {
					sqlBuffer.WriteString("`" + f + "`, ")
				}
			}
		}
		sqlBuffer.WriteString("FROM `" + sm.builder.model.tableName + "` ")
		//builder.wheres
		if ws := sm.wheresStr(); ws != "" {
			sqlBuffer.WriteString(ws)
		}
		//builder.orders
		if len(sm.builder.orders) > 0 {
			for i, item := range sm.builder.orders {
				if i == 0 {
					sqlBuffer.WriteString(fmt.Sprintf("ORDER BY `%s` %s ", item[1], item[0]))
				} else {
					sqlBuffer.WriteString(fmt.Sprintf(", `%s` %s ", item[1], item[0]))
				}
			}
		}
		//limit
		if sm.builder.limit > 0 {
			if sm.builder.limitOffset == 0 {
				sqlBuffer.WriteString("LIMIT 1 ")
			} else {
				sqlBuffer.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d ", sm.builder.limitOffset, sm.builder.limitOffset*(sm.builder.limit-1)))
			}
		}
	case OPUpdate:
		if len(sm.builder.updateFields) == 0 {
			return fmt.Errorf("edb StmtMysql.Build err: OPUpdate no updated fields")
		}
		sqlBuffer.WriteString(fmt.Sprintf("UPDATE `%s` SET ", sm.builder.model.tableName))

		updateStr := ""
		for _, item := range sm.builder.updateFields {
			if f, ok := sm.builder.model.entityFields[item]; ok {
				updateStr += ",`" + item + "` = ?"
				if f.fType == "time.Time" || f.fTagType == "Time" {
					if f.value == nil {
						continue
					}
					switch f.fTagType {
					case TagDate:
						sm.bindings = append(sm.bindings, f.value.(time.Time).Format(FTimeDate))
					case TagTime:
						sm.bindings = append(sm.bindings, f.value.(time.Time).Format(FTimeTime))
					case TagDateTime:
						fallthrough
					default:
						sm.bindings = append(sm.bindings, f.value.(time.Time).Format(FTimeDateTime))
					}
				} else {
					sm.bindings = append(sm.bindings, f.value)
				}

			}
		}
		sqlBuffer.WriteString(strings.TrimLeft(updateStr, ",") + " ")

		if len(sm.builder.wheres) == 0 {
			//use the pk as where condition
			pk := sm.builder.model.pkField
			if pk == "" {
				return fmt.Errorf("edb StmtMysql.Build err: OPUpdate the pirmary key cannot be fount")
			}
			sqlBuffer.WriteString(fmt.Sprintf("WHERE `%s` = ? ", pk))
			sm.bindings = append(sm.bindings, sm.builder.model.entityFields[pk].value)
		} else {
			//builder.wheres
			if ws := sm.wheresStr(); ws != "" {
				sqlBuffer.WriteString(ws)
			}
		}
	case OPInsert:
		sqlBuffer.WriteString("INSERT INTO `" + sm.builder.model.tableName + "` ")

		count := 0
		fstr := ""
		for _, f := range sm.builder.model.entityFields {
			if f.isAuto {
				continue
			}
			if f.fType == "time.Time" || f.fTagType == "Time" {
				if f.value == nil {
					continue
				}
				fstr += ",`" + f.name + "`"
				switch f.fTagType {
				case TagDate:
					sm.bindings = append(sm.bindings, f.value.(time.Time).Format(FTimeDate))
				case TagTime:
					sm.bindings = append(sm.bindings, f.value.(time.Time).Format(FTimeTime))
				case TagDateTime:
					fallthrough
				default:
					sm.bindings = append(sm.bindings, f.value.(time.Time).Format(FTimeDateTime))
				}
			} else {
				fstr += ",`" + f.name + "`"
				sm.bindings = append(sm.bindings, f.value)
			}
			count++
		}

		vstr := ""
		for i := 0; i < count; i++ {
			vstr += ",?"
		}

		sqlBuffer.WriteString(fmt.Sprintf("(%s) VALUES (%s)", strings.TrimLeft(fstr, ","), strings.TrimLeft(vstr, ",")))

	case OPDelete:
		sqlBuffer.WriteString(fmt.Sprintf("DELETE FROM `%s` ", sm.builder.model.tableName))
		if len(sm.builder.wheres) == 0 {
			//use the pk as where condition
			pk := sm.builder.model.pkField
			if pk != "" {
				sqlBuffer.WriteString(fmt.Sprintf("WHERE `%s` = ? ", pk))
				sm.bindings = append(sm.bindings, sm.builder.model.entityFields[pk].value)
			}
		} else {
			//builder.wheres
			if ws := sm.wheresStr(); ws != "" {
				sqlBuffer.WriteString(ws)
			}
		}
	default:
		return fmt.Errorf("edb StmtMysql.Build err: undefined OP type")
	}
	sqlBuffer.WriteString(";")
	sm.prepareSQL = sqlBuffer.String()
	return nil
}

// PrepareSQL get prepare sql
func (sm *StmtMysql) PrepareSQL() string {
	if sm.prepareSQL != "" {
		return sm.prepareSQL
	}

	return ""
}

// Bindings get prepare bind parameters
func (sm *StmtMysql) Bindings() []interface{} {
	return sm.bindings
}

func (sm *StmtMysql) wheresStr() string {
	l := len(sm.builder.wheres)
	if l == 0 {
		return ""
	}
	sql := "WHERE "
	s := make([]string, l)
	for i, w := range sm.builder.wheres {
		s[i] = fmt.Sprintf("`%s` %s ? ", w[0], w[1])
		sm.bindings = append(sm.bindings, w[2])
	}
	sql = sql + strings.Join(s, "AND ")
	return sql
}

func (sm *StmtMysql) reset() {
	sm.prepareSQL = ""
	sm.bindings = make([]interface{}, 0)
	sm.op = 0
}
