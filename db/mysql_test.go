package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func createDB() *MySQL {
	conf := Conf{
		IP:          "127.0.0.1",
		DB:          "gproxy",
		Port:        3306,
		UserName:    "gproxy",
		Password:    "gproxy123Aa!",
		Idle:        1,
		Active:      2,
		IdleTimeout: 10,
	}

	return NewMySQL(conf)
}

type TestTbl struct {
	ID      int    `json:"id" gorm:"column:id"`
	Name    string `json:"name" gorm:"column:name"`
	Age     int    `json:"age" gorm:"column:age"`
	Group   int    `json:"group" gorm:"column:group"`
	Company string `json:"company" gorm:"column:company"`
}

type User struct {
	ID       int64  `json:"id" gorm:"primaryKey;column:id"`
	UserName string `json:"username" gorm:"column:username"`
	Password string `json:"-" gorm:"column:password;size:64"`
}

type Resource struct {
	ID     int64  `json:"id" gorm:"primaryKey;column:id"`
	NodeID int64  `json:"nodeid" gorm:"column:nodeid"`
	IP     string `json:"ip" gorm:"column:ip"`
	MASK   string `json:"mask" gorm:"column:mask"`
}

//func TestCreateTbl(t *testing.T) {
//	sql := createDB()
//	if err := sql.orm.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(&model.PrivilegeGroup{}, &model.Company{}, &model.User{}, &model.Resource{}, &model.Version{}); err != nil {
//		t.Error(err)
//	}
//}

func TestMySQL(t *testing.T) {
	assert := assert.New(t)
	sql := createDB()

	assert.Nil(sql.dropTbl(&TestTbl{}))
	assert.Nil(sql.createTbl(&TestTbl{}))

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
		t.Logf("op:%s, id:%d", c.OP, c.Obj.ID)
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
				Offset: 1,
				Limit:  10,
			}
			assert.Equal(c.Err, sql.GetPages(&objs, pages, ""))
		case "Delete":
			assert.Equal(c.Err, sql.Delete(&TestTbl{}, "id = ?", c.Obj.ID))
		}
	}
}

func TestMySQLUpdate(t *testing.T) {
	assert := assert.New(t)
	sql := createDB()

	assert.Nil(sql.dropTbl(&TestTbl{}))
	assert.Nil(sql.createTbl(&TestTbl{}))

	tbl := TestTbl{
		Name: "test1",
		Age:  3,
	}

	assert.NoError(sql.Add(&tbl))
	assert.NoError(sql.Update(&tbl, "Age", 4, "id = ?", tbl.ID))
}

func TestMySQLMultiUpdate(t *testing.T) {
	assert := assert.New(t)
	sql := createDB()

	assert.Nil(sql.dropTbl(&TestTbl{}))
	assert.Nil(sql.createTbl(&TestTbl{}))

	tbl := TestTbl{
		Name:    "test1",
		Age:     3,
		Group:   4,
		Company: "this company",
	}

	assert.NoError(sql.Add(&tbl))

	newTbl := TestTbl{
		Group:   5,
		Company: "this is a new company",
	}
	assert.NoError(sql.Updates(&newTbl, newTbl, "name = ? and Age = 3", "test1", 3))
}

func TestMySQLGet(t *testing.T) {
	assert := assert.New(t)
	sql := createDB()

	assert.Nil(sql.dropTbl(&TestTbl{}))
	assert.Nil(sql.createTbl(&TestTbl{}))

	tbl := TestTbl{
		Name:    "test1",
		Age:     3,
		Group:   4,
		Company: "this company",
	}

	assert.NoError(sql.Add(&tbl))

	assert.NoError(sql.GetOne(&tbl, "name = ? and age = ?", tbl.Name, tbl.Age))
}

type TestUser struct {
	gorm.Model
	Name   string
	Groups []TestGroup `gorm:"many2many:user_group;"`
}

type TestGroup struct {
	gorm.Model
	Name      string
	TestUsers []TestUser `gorm:"many2many:user_group;"`
}

func TestHasMany(t *testing.T) {
	assert := assert.New(t)
	sql := createDB()

	sql.dropTbl(&TestUser{})
	sql.dropTbl(&TestGroup{})

	assert.NoError(sql.orm.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(&TestUser{}, &TestGroup{}))

	groups := []TestGroup{
		{
			Name: "Admin",
		},
		{
			Name: "Normal",
		},
	}

	u := TestUser{
		Name:   "admin",
		Groups: groups,
	}

	sql.Add(&u)

	// update groups
	groups = append(groups, TestGroup{
		Name: "Super",
	})
	u.Groups = groups
	assert.NoError(sql.UpdatesAll(&u))

	u1 := TestUser{}

	assert.NoError(sql.orm.Preload(clause.Associations).Find(&u1, "id = ?", u.ID).Error)
	assert.Equal(u.Groups, u1.Groups)
}
