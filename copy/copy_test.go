package copy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	assert := assert.New(t)
	type A struct {
		ID   int
		Name string
	}
	type B struct {
		A
		Addr string
	}

	v1 := B{
		A{
			ID:   3,
			Name: "name",
		},
		"addr1",
	}
	var v2 B

	assert.NoError(SameType(&v2, v1))
	assert.Equal(v1, v2)
}
