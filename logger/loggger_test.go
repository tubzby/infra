package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog2File(t *testing.T) {
	assert := assert.New(t)
	file := "/tmp/log.file"

	maxSize := 1024 * 1024
	rotate := 3
	assert.NoError(LogToFile(file, maxSize, uint(rotate)))

	bs := make([]byte, 1024*100)
	for i := 0; i < 100; i++ {
		Infof("str: %s", bs)
	}
}
