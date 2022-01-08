package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBits(t *testing.T) {
	assert := assert.New(t)

	b := 0b00010001

	assert.Equal(true, BitOn(b, 4))
	assert.Equal(false, BitOn(b, 5))

	SetBit(&b, 2)
	assert.Equal(0b00010101, b)

	ClearBit(&b, 4)
	assert.Equal(0b00000101, b)
}
