package edb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	//_ test
	_ "github.com/go-sql-driver/mysql"
)

//naming rules:
// database : user_config   => struct : UserConfig
//
// database table
// CREATE TABLE `user` (
// 	`id` int NOT NULL AUTO_INCREMENT,
// 	`name` varchar(50) DEFAULT '',
// 	`age` int unsigned DEFAULT '0',
// 	`created_at` datetime DEFAULT NULL,
// 	`updated_at` datetime DEFAULT NULL,
// 	PRIMARY KEY (`id`)
//  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

// to struct
// type User struct {
// 	Id        int `type:"autoPk"`
// 	Name      string
//  Age       int
// 	CreatedAt time.Time `type:"dateTime"`
// 	UpdatedAt time.Time `type:"dateTime"`
// }

func TestModelNew(t *testing.T) {
	TestBoot(t)

	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"dateTime"`
		UpdatedAt time.Time `type:"dateTime"`
	}

	u1 := &User{
		Name:      "tom",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m1, err := New(u1)
	assert.Nil(t, err)
	assert.True(t, m1.isAuto)
	assert.Equal(t, "id", m1.pkField)

	//uncorrect
	type (
		//User2
		User2 struct {
			Id  int `type:"autoPk"`
			Id2 int `type:"pk"`
		}
		//User3
		User3 struct {
			Id    int `type:"autoPk"`
			Names []string
		}
		//User4
		User4 struct {
			Id   int `type:"autoPk"`
			name string
		}
	)
	tests := []struct {
		entity    interface{}
		expectErr string
	}{
		{
			User{},
			"edb Model.checkEntity err: the parameter \"entity\" must be a struct pointer",
		},
		{
			&User2{},
			"edb Model.setTableAttributes err: has been set pk: id, can no longer set the filed `id2` as pk",
		},
		{
			&User3{},
			"edb Model.setTableAttributes err: unsupported field type: []string",
		},
		{
			&User4{},
			"edb Model.setTableAttributes err: field: name unexported",
		},
	}

	for _, item := range tests {
		_, err := New(item.entity)
		assert.EqualError(t, err, item.expectErr)
	}
}

func TestModelInsert(t *testing.T) {
	TestBoot(t)

	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"dateTime"`
		UpdatedAt time.Time `type:"dateTime"`
	}

	u := &User{
		Name:      "tom",
		Age:       222,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m, err := New(u)
	assert.Nil(t, err)

	m.Exec("truncate `user`;")

	id, err := m.Insert()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id)
}

func TestModelFirst(t *testing.T) {
	TestBoot(t)

	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"dateTime"`
		UpdatedAt time.Time `type:"dateTime"`
	}

	u1 := &User{}
	m1, err := New(u1)
	assert.Nil(t, err)

	m1.Exec("truncate `user`;")

	entity, err := m1.First()
	assert.Nil(t, err)
	assert.Nil(t, entity)

	u2 := &User{
		Name:      "tom",
		Age:       222,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m2, err := New(u2)
	m2.Insert()

	m1.Eq("name", "tom")
	e2, err := m1.First()
	assert.Nil(t, err)
	assert.NotNil(t, e2)
	assert.Equal(t, "tom", e2.(*User).Name)
	assert.Equal(t, 222, e2.(*User).Age)
}

func TestModelGet(t *testing.T) {
	TestBoot(t)

	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"dateTime"`
		UpdatedAt time.Time `type:"dateTime"`
	}

	u1 := &User{}
	m1, err := New(u1)
	assert.Nil(t, err)

	m1.Exec("truncate `user`;")

	collect, err := m1.Get()
	assert.Nil(t, err)
	assert.False(t, collect.Next())

	tt, _ := time.ParseInLocation(FTimeDateTime, "2021-01-01 01:01:01", time.Local)
	u2 := &User{
		Name:      "tom",
		Age:       222,
		CreatedAt: tt,
		UpdatedAt: tt,
	}
	m2, err := New(u2)
	id, err := m2.Insert()
	assert.Equal(t, int64(1), id)
	id2, err := m2.Insert()
	assert.Equal(t, int64(2), id2)

	collect2, err := m1.Get()
	assert.Nil(t, err)

	count := 0
	ids := []int{0, 0}
	for collect2.Next() {
		item := collect2.Item()
		assert.NotNil(t, item)
		ids[count] = item.(*User).Id
		count++
	}
	assert.Equal(t, 2, count)
	assert.Equal(t, 1, ids[0])
	assert.Equal(t, 2, ids[1])

}

func TestModelUpdate(t *testing.T) {

	TestBoot(t)

	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"dateTime"`
		UpdatedAt time.Time `type:"dateTime"`
	}

	u1 := &User{
		Name: "tom",
		Age:  222,
	}
	m1, err := New(u1)
	assert.Nil(t, err)

	m1.Exec("truncate `user`;")
	id, err := m1.Insert()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id)

	m2, err := New(&User{
		Id:   1,
		Name: "ttt",
		Age:  333,
	})
	_, err2 := m2.Update([]string{})
	assert.EqualError(t, err2, "edb StmtMysql.Build err: OPUpdate no updated fields")
	//default  pk    => where `id` = 1
	rowAffected, err := m2.Update([]string{"name", "age"})
	assert.Equal(t, int64(1), rowAffected)

	e, err := m1.Eq("id", 1).First()
	assert.Nil(t, err)
	assert.Equal(t, "ttt", e.(*User).Name)

	m3, err := New(&User{
		Name: "tom",
		Age:  444,
	})
	//specify conditions => where name = 'ttt'
	m3.Eq("name", "ttt").Update([]string{"name", "age"})

	e2, err := m1.Eq("id", 1).First()
	assert.Nil(t, err)
	assert.Equal(t, "tom", e2.(*User).Name)
	assert.Equal(t, 444, e2.(*User).Age)
}

func TestModelDelete(t *testing.T) {

	TestBoot(t)
	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"dateTime"`
		UpdatedAt time.Time `type:"dateTime"`
	}
	m1, err := New(&User{
		Id: 1,
	})
	assert.Nil(t, err)
	m1.Exec("truncate `user`;")

	m2, err := New(&User{
		Name: "ttt",
		Age:  333,
	})
	id, err := m2.Insert()
	assert.Equal(t, int64(1), id)
	m3, err := New(&User{
		Name: "tt",
		Age:  22,
	})
	id2, err := m3.Insert()
	assert.Equal(t, int64(2), id2)

	//default  pk    => where `id` = 1
	rowAffected, err := m1.Delete()
	assert.Equal(t, int64(1), rowAffected)

	//specify conditions => where age = 22
	rowAffected2, err := m1.Eq("age", 22).Delete()
	assert.Equal(t, int64(1), rowAffected2)

	e, err := m1.First()
	assert.Nil(t, e)
}

func TestModelPaginate(t *testing.T) {
	TestBoot(t)
	type User struct {
		Id        int `type:"autoPk"`
		Name      string
		Age       int
		CreatedAt time.Time `type:"dateTime"`
		UpdatedAt time.Time `type:"dateTime"`
	}

	m1, err := New(&User{})
	assert.Nil(t, err)
	m1.Exec("truncate `user`;")

	for i := 0; i < 28; i++ {
		m, err := New(&User{
			Name: "tom",
			Age:  i,
		})
		assert.Nil(t, err)
		m.Insert()
	}

	c, err := m1.Gte("age", 10).Paginate(1, 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(18), c.Total())
	pageCount := 0
	for c.Next() {
		pageCount++
	}
	assert.Equal(t, 10, pageCount)

	c2, err := m1.Gte("age", 10).Paginate(2, 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(18), c.Total())
	pageCount2 := 0
	for c2.Next() {
		pageCount2++
	}
	assert.Equal(t, 8, pageCount2)

}
