package db

// Connecter interface
type Connecter interface {
	Add(obj interface{}) error
	GetOne(obj interface{}, query string, args ...interface{}) error
	Delete(obj interface{}, query string, args ...interface{}) error
	GetPages(objs interface{}, query Query) error
	GetAll(objs interface{}, query string, args ...interface{}) error
	Count(obj interface{}, query Query) (int64, error)
	// update single column
	Update(obj interface{}, column string, value interface{}, query string, args ...interface{}) error
	// update multiple columns
	Updates(obj interface{}, values interface{}, query string, args ...interface{}) error
	UpdatesAll(obj interface{}) error
}

// CustomObj .
type CustomObj interface {
	TblName() string
}

type Query struct {
	Offset int
	Limit  int
	Filter map[string]interface{}
}
