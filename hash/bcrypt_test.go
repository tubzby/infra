package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcrypt(t *testing.T) {
	assert := assert.New(t)
	h := bcryptHasher{}
	password := []byte("123456")

	out, err := h.Hash(password)
	assert.NoError(err)
	assert.True(h.Match(out, password))

	t.Logf("hashed: %s", out)
}
