package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheMem(t *testing.T) {
	assert := assert.New(t)
	m := NewMemCache()

	type A struct {
		ID   int
		Name string
		Addr string
	}
	key := "key"
	value := A{
		2,
		"name",
		"Addr",
	}

	cases := []struct {
		OP    string
		Key   string
		Value A
		Err   error
		OK    bool
	}{
		{
			OP:    "Exist",
			Key:   key,
			Value: value,
			OK:    false,
		},
		{
			OP:    "Load",
			Key:   key,
			Value: value,
			Err:   ErrNil,
		},
		{
			OP:    "Save",
			Key:   key,
			Value: value,
			Err:   nil,
		},
		{
			OP:    "Exist",
			Key:   key,
			Value: value,
			OK:    true,
		},
		{
			OP:    "Load",
			Key:   key,
			Value: value,
			Err:   nil,
		},
		{
			OP:  "Delete",
			Key: key,
			Err: nil,
		},
		{
			OP:  "Exist",
			Key: key,
			OK:  false,
		},
		{
			OP:    "Load",
			Key:   key,
			Value: value,
			Err:   ErrNil,
		},
	}

	for _, c := range cases {
		switch c.OP {
		case "Exist":
			assert.Equal(c.OK, m.Exist(c.Key))
		case "Load":
			var obj A
			assert.Equal(c.Err, m.Load(c.Key, &obj))
			if c.Err == nil {
				assert.Equal(c.Value, obj)
			}
		case "Save":
			assert.Equal(c.Err, m.Save(c.Key, c.Value, 0))
		case "Delete":
			assert.Equal(c.Err, m.Delete(c.Key))
		default:
			t.Errorf("invalid op:%s", c.OP)
		}
	}
}
