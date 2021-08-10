package edb

type (

	// Builder builder
	Builder struct {
		model        *Model
		fields       []string
		wheres       [][]interface{}
		updateFields []string
		orders       [][]string
		limit        int64
		limitOffset  int64
		tableName    string
	}
)

// NewBuilder new builder
func NewBuilder() *Builder {
	return &Builder{
		fields:       make([]string, 0),
		wheres:       make([][]interface{}, 0),
		orders:       make([][]string, 0),
		updateFields: make([]string, 0),
		limitOffset:  10,
	}
}

// reset reset builder attr
func (b *Builder) reset() {
	b.fields = make([]string, 0)
	b.wheres = make([][]interface{}, 0)
	b.updateFields = make([]string, 0)
	b.orders = make([][]string, 0)
	b.limit = 0
	b.limitOffset = 10
}

// Select select field
func (b *Builder) Select(fields []string) error {
	b.fields = fields
	return nil
}

// WhereCondition where condition
func (b *Builder) WhereCondition(field string, condition string, value interface{}) error {
	b.wheres = append(b.wheres, []interface{}{field, condition, value})
	return nil
}

// OrderBy ASC sort
func (b *Builder) OrderBy(field string) {
	b.orders = append(b.orders, []string{"ASC", field})
}

// OrderByDesc DESC sort
func (b *Builder) OrderByDesc(field string) {
	b.orders = append(b.orders, []string{"DESC", field})
}

// Update update opreate, pass the fields that need to be updated
func (b *Builder) Update(updateFields []string) {
	b.updateFields = updateFields
}
