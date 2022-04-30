package cache

import (
	"gitee.com/romeo_zpl/infra/logger"
	"github.com/vmihailenco/msgpack"
)

// Mem in mem cache implement cacher
type Mem struct {
	mem map[string][]byte
}

// NewMemCache Create MemCache
func NewMemCache() *Mem {
	return &Mem{
		mem: make(map[string][]byte),
	}
}

// Save key value
func (cm *Mem) Save(key string, obj interface{}, expire int) error {
	logger.Debugf("Save key:%s", key)
	// load to redis
	bs, err := msgpack.Marshal(obj)
	if err != nil {
		logger.Errorf("msgpack failed(%v)", err)
		return ErrMarshal
	}

	cm.mem[key] = bs
	return nil
}

// Load obj with key
func (cm *Mem) Load(key string, obj interface{}) error {
	bs, ok := cm.mem[key]
	if !ok {
		return ErrNil
	}

	if err := msgpack.Unmarshal(bs, obj); err != nil {
		logger.Errorf("unpack error(%v)", err)
		return ErrUnMarshal
	}
	return nil
}

// Delete obj with key
func (cm *Mem) Delete(key string) error {
	logger.Debugf("Save key:%s", key)
	delete(cm.mem, key)
	return nil
}

// Exist check if key exist
func (cm *Mem) Exist(key string) bool {
	_, ok := cm.mem[key]
	return ok
}
