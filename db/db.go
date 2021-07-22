package db

// Connecter interface
type Connecter interface {
	Add(obj interface{}) error
	GetOne(obj interface{}, query string, args ...interface{}) error
	Delete(obj interface{}, query string, args ...interface{}) error
	GetPages(objs interface{}, page PageParam, query string, args ...interface{}) error
	GetAll(objs interface{}, query string, args ...interface{}) error
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

// PageParam is parameter for querying a list of resource
type PageParam struct {
	PageNo   int `json:"pageno" form:"pageno"`
	PageSize int `json:"pagesize" form:"pagesize"`
}
