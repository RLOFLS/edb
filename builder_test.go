package edb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	builder := NewBuilder()
	builder.Select([]string{"id", "name"})
	assert.Equal(t, 2, len(builder.fields))

	builder.Select([]string{"id"})
	assert.Equal(t, 1, len(builder.fields))
}

func TestWhereCondition(t *testing.T) {
	builder := NewBuilder()
	builder.WhereCondition("id", "=", 1)
	assert.Equal(t, 1, len(builder.wheres))

	builder.WhereCondition("name", "LIKE", "%t")
	assert.Equal(t, 2, len(builder.wheres))
}

func TestOrder(t *testing.T) {
	builder := NewBuilder()
	builder.OrderBy("id")
	builder.OrderBy("name")
	builder.OrderByDesc("age")

	e := [][]string{
		{"ASC", "id"},
		{"ASC", "name"},
		{"DESC", "age"},
	}
	for i, item := range builder.orders {
		assert.Equal(t, e[i][0], item[0])
		assert.Equal(t, e[i][1], item[1])
	}
}

func TestUpdate(t *testing.T) {
	b := NewBuilder()
	b.Update([]string{"id", "name"})
	assert.Equal(t, []string{"id", "name"}, b.updateFields)

	b.Update([]string{"age"})
	assert.Equal(t, []string{"age"}, b.updateFields)
}
