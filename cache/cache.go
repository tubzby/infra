package cache

// Cacher interface
type Cacher interface {
	Save(key string, obj interface{}, expire int) error
	Load(key string, obj interface{}) error
	Delete(key string) error
	Exist(key string) bool
}
