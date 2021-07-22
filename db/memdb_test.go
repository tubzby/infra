package db

import (
	"reflect"
	"testing"

	"gdcx.com/infra/copy"
	"github.com/stretchr/testify/assert"
)

func TestNewObj(t *testing.T) {
	assert := assert.New(t)
	type A struct {
		ID   int
		Name string
	}

	a := A{
		ID:   3,
		Name: "348",
	}

	func(obj interface{}) {
		newObj := reflect.New(reflect.TypeOf(obj).Elem())
		assert.Nil(copy.SameType(newObj.Interface(), obj))
		assert.Equal(newObj.Interface(), obj)
	}(&a)

	func(obj interface{}) {
		val := reflect.Indirect(reflect.ValueOf(obj))
		assert.Equal(3, val.FieldByName("ID").Interface())
	}(&a)
}

func TestCopyStruct(t *testing.T) {
	assert := assert.New(t)
	type A struct {
		ID   int
		Name string
	}

	a := A{
		ID:   3,
		Name: "348",
	}

	var out interface{}
	func(obj interface{}) {
		out = reflect.Indirect(reflect.ValueOf(obj)).Interface()
	}(&a)

	assert.Equal(out, a)

	a.ID = 4
	assert.NotEqual(out, a)
}

func TestGetFieldByTag(t *testing.T) {
	assert := assert.New(t)
	type A struct {
		ID   int    `json:"id" gorm:"primaryKey;column:id"`
		Name string `json:"name" gorm:"column:columnname"`
	}

	a := A{
		ID:   3,
		Name: "name",
	}

	assert.Equal(0, getField(a, "ID"))
	assert.Equal(0, getField(a, "id"))
	assert.Equal(1, getField(a, "Name"))
	assert.Equal(1, getField(a, "columnname"))
	assert.Equal(-1, getField(a, "idd"))
}

func TestNestedStruct(t *testing.T) {
	type A struct {
		ID   int    `json:"id" gorm:"autoIncrement;primaryKey;column:id"`
		Name string `json:"name" gorm:"column:name"`
	}

	type B struct {
		A
		Addr string `json:"addr" gorm:"column:addr"`
	}

	assert := assert.New(t)
	sql := NewMemDB()

	a := B{
		A{
			ID:   2,
			Name: "name",
		},
		"address",
	}

	assert.NoError(sql.Add(&a))

	var r B
	assert.NoError(sql.GetOne(&r, "id = ?", 2))
	assert.Equal(a.Name, r.Name)
}

func TestGormTag(t *testing.T) {
	//assert := assert.New(t)
	type A struct {
		ID   int    `json:"id" gorm:"column:id"`
		Name string `json:"name" gorm:"column:name"`
	}

	a := A{
		ID:   3,
		Name: "348",
	}

	func(obj interface{}) {
		typ := reflect.TypeOf(obj).Elem()
		for i := 0; i < typ.NumField(); i++ {
			typ.Field(i).Tag.Get("gorm")
		}
	}(&a)

}
func TestMemDB(t *testing.T) {
	assert := assert.New(t)
	sql := NewMemDB()

	cases := []struct {
		OP  string
		Obj TestTbl
		Err error
	}{
		{
			OP: "GetOne",
			Obj: TestTbl{
				ID: 1,
			},
			Err: ErrNil,
		},
		{
			OP: "Add",
			Obj: TestTbl{
				ID:   1,
				Name: "zp",
				Age:  34,
			},
			Err: nil,
		},
		{
			OP: "Add",
			Obj: TestTbl{
				ID:   2,
				Name: "shirley",
				Age:  34,
			},
			Err: nil,
		},
		{
			OP: "GetOne",
			Obj: TestTbl{
				ID:   2,
				Name: "shirley",
				Age:  34,
			},
			Err: nil,
		},
		{
			OP: "GetOne",
			Obj: TestTbl{
				ID:   1,
				Name: "zp",
				Age:  34,
			},
			Err: nil,
		},
		{
			OP:  "GetPages",
			Err: nil,
		},
		{
			OP: "Delete",
			Obj: TestTbl{
				ID: 2,
			},
			Err: nil,
		},
		{
			OP: "GetOne",
			Obj: TestTbl{
				ID: 2,
			},
			Err: ErrNil,
		},
	}

	for _, c := range cases {
		switch c.OP {
		case "GetOne":
			var t TestTbl
			assert.Equal(c.Err, sql.GetOne(&t, "id = ?", c.Obj.ID))
			if c.Err == nil {
				assert.Equal(c.Obj.Name, t.Name)
				assert.Equal(c.Obj.Age, t.Age)
			}
		case "Add":
			assert.Equal(c.Err, sql.Add(&c.Obj))
		case "GetPages":
			var objs []TestTbl
			pages := PageParam{
				PageNo:   1,
				PageSize: 10,
			}
			assert.Equal(c.Err, sql.GetPages(&objs, pages, ""))
		case "Delete":
			assert.Equal(c.Err, sql.Delete(&TestTbl{}, "id = ?", c.Obj.ID))
		}
	}
}
