package logger

import "github.com/sirupsen/logrus"

type Entry struct {
	e *logrus.Entry
}

func WithFields(fields map[string]interface{}) *Entry {
	return &Entry{
		e: logrus.WithFields(fields),
	}
}

// Error .
func (e *Entry) Error(v ...interface{}) {
	e.e.Error(v...)
}

// Errorf .
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.e.Errorf(format, args...)
}

// Fatal .
func (e *Entry) Fatal(v ...interface{}) {
	logrus.Panic(v...)
}

// Fatalf .
func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.e.Fatalf(format, args...)
}

// Info .
func (e *Entry) Info(v ...interface{}) {
	e.e.Info(v...)
}

// Infof .
func (e *Entry) Infof(format string, args ...interface{}) {
	e.e.Infof(format, args...)
}

// Debug .
func (e *Entry) Debug(v ...interface{}) {
	e.e.Debug(v...)
}

// Debugf .
func (e *Entry) Debugf(format string, args ...interface{}) {
	e.e.Debugf(format, args...)
}
