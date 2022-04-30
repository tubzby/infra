package cache

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var conf Conf = Conf{
	IP:       "192.168.3.100",
	Port:     6379,
	Password: "zplredis",
}

func TestCache(t *testing.T) {
	assert := assert.New(t)
	k := "testKey"
	v := []struct {
		ID   int
		Name string
	}{
		{
			ID:   3,
			Name: "this is zpl speaking...",
		},
		{
			ID:   0,
			Name: "",
		},
	}

	r := NewRedis(conf)
	assert.NoError(r.Save(k, v[0], 100))
	assert.True(r.Exist(k))
	assert.NoError(r.Load(k, &v[1]))
	assert.True(reflect.DeepEqual(v[0], v[1]))

	assert.NoError(r.Delete(k))
	assert.False(r.Exist(k))
	assert.Error(r.Load(k, &v[1]))
}
