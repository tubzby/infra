package cache

import (
	"fmt"
	"time"

	"gitee.com/romeo_zpl/infra/logger"
	"github.com/go-redis/redis"
	"github.com/vmihailenco/msgpack"
)

// Conf is redis configurations
type Conf struct {
	IP       string
	Port     int
	Password string
	DB       int
}

// Redis .
type Redis struct {
	redis *redis.Client
}

// NewRedis create cacher
func NewRedis(conf Conf) *Redis {
	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.IP, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})
	return &Redis{
		redis: r,
	}
}

// Save to redis
func (r *Redis) Save(key string, obj interface{}, expire int) error {
	// load to redis
	bs, err := msgpack.Marshal(obj)
	if err != nil {
		logger.Errorf("msgpack failed(%v)", err)
		return ErrMarshal
	}
	if err = r.redis.Set(key, bs, time.Duration(expire)*time.Second).Err(); err != nil {
		logger.Errorf("add category to redis failed(%v)", err)
		return err
	}
	return nil
}

// Load from redis
func (r *Redis) Load(key string, obj interface{}) error {
	bs, err := r.redis.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrNil
		}
		logger.Errorf("load key(%s) from redis error(%v)", key, err)
		return err
	}

	if err := msgpack.Unmarshal(bs, obj); err != nil {
		logger.Errorf("unpack error(%v)", err)
		return ErrUnMarshal
	}
	return nil
}

// Exist check if key exist
func (r *Redis) Exist(key string) bool {
	ret, err := r.redis.Exists(key).Result()
	if err != nil {
		logger.Errorf("redis exist error(%v)", err)
		return false
	}
	return ret == 1
}

// Delete key
func (r *Redis) Delete(key string) error {
	if err := r.redis.Del(key).Err(); err != nil {
		logger.Errorf("redis key(%s) del error(%v)", key, err)
		return err
	}
	return nil
}
