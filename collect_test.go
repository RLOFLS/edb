package edb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	//_ test
	_ "github.com/go-sql-driver/mysql"
)

// In addition to regular operations like model => model_test.go file
// you can alse customize structure mapping
// for test, just QueryCollect
func TestMappingStructure(t *testing.T) {
	TestBoot(t)

	type UserCount struct {
		UserName string
		Total    int
	}

	m, err := New(&UserCount{})
	assert.Nil(t, err)

	m.Exec("truncate `user`;")

	//insert data
	// name = 'tom' => 8 records
	// name = 'tom2' => 5 records
	testInsertDataForMappingStructure(t)

	prepareSQL := "SELECT `name` AS user_name, count(*) AS total FROM `user` GROUP BY `name`"
	c, err := m.QueryCollect(prepareSQL)
	assert.Nil(t, err)

	assert.True(t, c.Next())

	i := c.Item()
	assert.Equal(t, "tom", i.(*UserCount).UserName)
	assert.Equal(t, 8, i.(*UserCount).Total)

	assert.True(t, c.Next())
	i2 := c.Item()
	assert.Equal(t, "tom2", i2.(*UserCount).UserName)
	assert.Equal(t, 5, i2.(*UserCount).Total)

	assert.False(t, c.Next())

}

func testInsertDataForMappingStructure(t *testing.T) {

	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"dateTime"`
		UpdatedAt time.Time `type:"dateTime"`
	}

	//insert data
	for i := 0; i < 8; i++ {
		m, err := New(&User{
			Name: "tom",
		})
		assert.Nil(t, err)
		m.Insert()
	}
	for i := 0; i < 5; i++ {
		m, err := New(&User{
			Name: "tom2",
		})
		assert.Nil(t, err)
		m.Insert()
	}
}
