package db

import (
	"errors"
	"reflect"
	"strings"

	"gitee.com/romeo_zpl/infra/copy"
	"gitee.com/romeo_zpl/infra/logger"
)

// TODO
// support gorm fieldname
// reduce copy of mem

var _ Connecter = new(MemDB)

// MemDB is db in memory
type MemDB struct {
	mem map[string][]interface{}
}

// NewMemDB create memory db
func NewMemDB() *MemDB {
	return &MemDB{
		mem: make(map[string][]interface{}),
	}
}

// Add obj to db
func (mdb *MemDB) Add(obj interface{}) error {
	if !checkType(obj) {
		return ErrParam
	}

	tbl := mdb.getTblName(obj)

	_, ok := mdb.mem[tbl]
	if !ok {
		mdb.mem[tbl] = make([]interface{}, 0)
	}
	mdb.mem[tbl] = append(mdb.mem[tbl], reflect.Indirect(reflect.ValueOf(obj)).Interface())

	return nil
}

// GetOne from db
// current only support schema as Where("field1 = ?", v1)
func (mdb *MemDB) GetOne(obj interface{}, query string, args ...interface{}) error {
	tbl := mdb.getTblName(obj)
	v, ok := mdb.mem[tbl]
	if !ok {
		return ErrNil
	}

	found := false
	fieldName := strings.Fields(query)[0]
	mdb.matchField(v, fieldName, args[0], func(idx int, o interface{}) {
		if err := copy.SameType(obj, o); err != nil {
			logger.Errorf("copy error(%v)", err)
		} else {
			found = true
		}
	})
	if !found {
		return ErrNil
	}

	return nil
}

// Delete from db
func (mdb *MemDB) Delete(obj interface{}, query string, args ...interface{}) error {
	tbl := mdb.getTblName(obj)
	v, ok := mdb.mem[tbl]
	if !ok {
		return ErrNil
	}

	fieldName := strings.Fields(query)[0]
	mdb.matchField(v, fieldName, args[0], func(idx int, o interface{}) {
		v[idx] = v[len(v)-1]
		mdb.mem[tbl] = v[:len(v)-1]
	})
	return nil
}

// GetPages from db
func (mdb *MemDB) GetPages(oobjs interface{}, query Query) error {
	return nil
}

func (mdb *MemDB) Count(obj interface{}, query Query) (int64, error) {
	return 0, nil
}

func (mdb *MemDB) GetAll(objs interface{}, query string, args ...interface{}) error {
	return nil
}

func checkType(obj interface{}) bool {
	return reflect.TypeOf(obj).Kind() == reflect.Ptr
}

func (mdb *MemDB) getTblName(obj interface{}) string {
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		return typ.Elem().Name() + "s"
	}
	return typ.Name() + "s"
}

func (mdb *MemDB) matchField(objs []interface{}, name string, value interface{}, f func(idx int, o interface{})) {
	for idx := range objs {
		// rvalue := reflect.Indirect(reflect.ValueOf(objs[idx]))
		// fieldIdx := getField(objs[idx], name)
		// if fieldIdx < 0 {
		// 	logger.Errorf("field %s not found", name)
		// 	continue
		// }
		// fieldValue := rvalue.Field(fieldIdx)
		// if value == fieldValue.Interface() {
		// 	f(idx, objs[idx])
		// }
		styp := reflect.TypeOf(objs[idx])
		svalue := reflect.Indirect(reflect.ValueOf(objs[idx]))
		walkObj(objs[idx], idx, styp, svalue, name, value, f)
	}
	return
}

func getField(obj interface{}, name string) int {
	typ := reflect.TypeOf(obj)
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Type.Kind() == reflect.Struct && f.Anonymous {

		}
		if name == f.Name || matchTag(name, f.Tag.Get("gorm")) {
			return i
		}
	}

	return -1
}

func walkObj(obj interface{}, idx int, subt reflect.Type, subv reflect.Value, name string, value interface{}, f func(idx int, o interface{})) {
	for i := 0; i < subt.NumField(); i++ {
		field := subt.Field(i)
		v := subv.Field(i)
		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			walkObj(obj, idx, field.Type, v, name, value, f)
		} else if name == field.Name || matchTag(name, field.Tag.Get("gorm")) {
			if v.Interface() == value {
				f(idx, obj)
			}
		}
	}
}

// match gorm name
func matchTag(name, tag string) bool {
	if len(tag) == 0 {
		return false
	}
	s1 := strings.Split(tag, ";")
	if len(s1) == 0 {
		return false
	}
	for _, s2 := range s1 {
		s3 := strings.Split(s2, ":")
		if len(s3) >= 2 && s3[0] == "column" && s3[1] == name {
			return true
		}
	}
	return false
}

func (mdb *MemDB) Update(obj interface{}, column string, value interface{}, query string, args ...interface{}) error {
	return errors.New("unimplemented")
}

func (mdb *MemDB) Updates(obj interface{}, values interface{}, query string, args ...interface{}) error {
	return errors.New("unimplemented")
}

func (mdb *MemDB) UpdatesAll(obj interface{}) error {
	return errors.New("unimplemented")
}
