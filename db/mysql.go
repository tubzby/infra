package db

import (
	"fmt"
	"reflect"
	"time"

	"gitee.com/romeo_zpl/infra/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Conf is db configurations
type Conf struct {
	IP          string
	Port        int
	UserName    string
	Password    string
	DB          string
	Idle        int
	Active      int
	IdleTimeout int
	Debug       bool
}

// MySQL .
type MySQL struct {
	conf    Conf
	orm     *gorm.DB
	preload bool
}

var _ Connecter = new(MySQL)

// NewMySQL create mysqldb
func NewMySQL(conf Conf) *MySQL {
	sql := MySQL{
		conf:    conf,
		preload: true,
	}
	sql.init()
	return &sql
}

func (sql *MySQL) init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		sql.conf.UserName,
		sql.conf.Password,
		sql.conf.IP,
		sql.conf.Port,
		sql.conf.DB)

	orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Errorf("db dsn(%s) err(%v)", dsn, err)
		return
	}

	mysql, _ := orm.DB()
	mysql.SetMaxIdleConns(sql.conf.Idle)
	mysql.SetMaxOpenConns(sql.conf.Active)
	mysql.SetConnMaxIdleTime(time.Duration(sql.conf.IdleTimeout) * time.Second)

	if sql.conf.Debug {
		logger.Info("db debug on")
		orm = orm.Debug()
	}
	sql.orm = orm
}

// Add add object to db
func (sql *MySQL) Add(obj interface{}) error {
	if !sql.checkState() {
		return ErrConnect
	}

	if err := sql.selectTbl(obj).Create(obj).Error; err != nil {
		logger.Errorf("add object(%v) error(%v)", obj, err)
		return err
	}
	return nil
}

// GetOne get a single object
func (sql *MySQL) GetOne(obj interface{}, query string, args ...interface{}) error {
	if !sql.checkState() {
		return ErrConnect
	}

	orm := sql.selectTbl(obj)
	if sql.preload {
		orm = orm.Preload(clause.Associations)
	}
	switch err := orm.Where(query, args...).First(obj).Error; err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return ErrNil
	default:
		logger.Errorf("sql error(%s)", err)
		return err
	}
}

// GetPages get a page of object
func (sql *MySQL) GetPages(objs interface{}, query Query) error {
	if !sql.checkState() {
		return ErrConnect
	}

	orm := sql.selectTbl(objs)
	if sql.preload {
		orm = orm.Preload(clause.Associations)
	}

	var err error
	if query.Filter != nil {
		orm = orm.Where(query.Filter)
	}

	if query.Limit > 0 {
		orm = orm.Limit(query.Limit).Offset(query.Offset)
	}
	err = orm.Find(objs).Error
	if err != nil {
		logger.Errorf("read object from db failed(%v)", err)
		return err
	}
	return nil
}

func (sql *MySQL) Count(obj interface{}, query Query) (int64, error) {
	if !sql.checkState() {
		return 0, ErrConnect
	}

	orm := sql.selectTbl(obj)
	if sql.preload {
		orm = orm.Preload(clause.Associations)
	}

	orm = orm.Model(obj)

	var count int64
	if query.Filter != nil {
		if db := orm.Where(query.Filter).Count(&count); db.Error != nil {
			return 0, db.Error
		}
	} else {
		if db := orm.Count(&count); db.Error != nil {
			return 0, db.Error
		}
	}
	return count, nil
}

// Delete one record
func (sql *MySQL) Delete(obj interface{}, query string, args ...interface{}) error {
	if !sql.checkState() {
		return ErrConnect
	}

	if len(query) == 0 {
		logger.Errorf("try to delete without conditions")
		return ErrParam
	}

	return sql.selectTbl(obj).Delete(obj, query, args).Error
}

func (sql *MySQL) GetAll(objs interface{}, query string, args ...interface{}) error {
	if !sql.checkState() {
		return ErrConnect
	}

	orm := sql.selectTbl(objs)

	if sql.preload {
		orm = orm.Preload(clause.Associations)
	}

	var err error
	if len(query) > 0 {
		err = orm.Where(query, args).Find(objs).Error
	} else {
		err = orm.Find(objs).Error
	}
	if err != nil {
		logger.Errorf("read obj from db failed(%s)", err)
		return err
	}
	return nil
}

func (sql *MySQL) DropTable(obj interface{}) error {
	if !sql.tblExist(obj) {
		return nil
	}

	if err := sql.orm.Migrator().DropTable(obj); err != nil {
		logger.Errorf("dropTbl error(%s)", err)
		return err
	}
	return nil
}

func (sql *MySQL) tblExist(obj interface{}) bool {
	return sql.orm.Migrator().HasTable(obj)
}
func (sql *MySQL) CreateTable(obj interface{}) error {
	if err := sql.orm.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(obj); err != nil {
		logger.Errorf("createTbl error(%s)", err)
		return err
	}
	return nil
}

func (sql *MySQL) selectTbl(obj interface{}) *gorm.DB {
	custom, ok := obj.(CustomObj)
	if ok {
		return sql.orm.Table(custom.TblName())
	} else {
		// for pointer to slice, detect if element is CutomObj type
		items := reflect.TypeOf(obj)
		if items.Kind() == reflect.Pointer {
			items = items.Elem()
		}
		kind := items.Kind()
		if kind != reflect.Slice && kind != reflect.Array {
			return sql.orm
		}

		// create a new object
		t := reflect.New(items.Elem())
		custom, ok = t.Interface().(CustomObj)
		if ok {
			return sql.orm.Table(custom.TblName())
		}
	}
	return sql.orm
}

func (sql *MySQL) Update(obj interface{}, column string, value interface{}, query string, args ...interface{}) error {
	if !sql.checkState() {
		return ErrConnect
	}

	orm := sql.selectTbl(obj)
	var err error

	if len(query) > 0 {
		err = orm.Model(obj).Where(query, args).Update(column, value).Error
	} else {
		err = orm.Model(obj).Where(obj).Update(column, value).Error
	}
	if err != nil {
		logger.Errorf("update error: %s", err)
	}
	return err
}

func (sql *MySQL) Updates(obj interface{}, values interface{}, query string, args ...interface{}) error {
	if !sql.checkState() {
		return ErrConnect
	}

	orm := sql.selectTbl(obj)

	return orm.Model(obj).Where(query, args...).Updates(values).Error
}

func (sql *MySQL) UpdatesAll(obj interface{}) error {
	if !sql.checkState() {
		return ErrConnect
	}

	orm := sql.selectTbl(obj)

	return orm.Save(obj).Error
}

func (sql *MySQL) valid() bool {
	return sql.orm != nil
}

func (sql *MySQL) checkState() bool {
	if !sql.valid() {
		sql.init()
	}
	return sql.valid()
}
