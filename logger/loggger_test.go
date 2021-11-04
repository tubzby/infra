package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
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

func TestLogTimeStamp(t *testing.T) {
	assert := assert.New(t)
	file := "/tmp/log.file"

	maxSize := 1024 * 1024
	rotate := 3
	assert.NoError(LogToFile(file, maxSize, uint(rotate)))
	Info("hello")
}

func TestLogWithField(t *testing.T) {
	dbg := logrus.New()
	dbg.SetLevel(logrus.DebugLevel)

	dbg.Debugf("hello there")

	mylog1 := dbg.WithFields(
		logrus.Fields{
			"id":  "134",
			"mod": "manager",
		},
	)
	mylog1.Infof("hello beautiful!")

	mylog2 := dbg.WithFields(
		logrus.Fields{
			"id":  "245",
			"mod": "modbus",
		},
	)
	mylog2.Infof("hello handsome!")
}

func TestEntry(t *testing.T) {
	EnableDebug()

	Debugf("A debug message")

	m1 := WithFields(map[string]interface{}{
		"id":  1234,
		"mod": "mod1",
	})

	m1.Debugf("mod1 debug message")

	m2 := WithFields(map[string]interface{}{
		"id":  4321,
		"mod": "mod2",
	})

	m2.Debugf("mod2 debug message")
}
