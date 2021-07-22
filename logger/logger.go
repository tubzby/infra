package logger

import (
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type myHook struct {
	FileName string
	Line     string
	Skip     int
	levels   []logrus.Level
}

//实现 logrus.Hook 接口
func (hook *myHook) Fire(entry *logrus.Entry) error {
	fileName, line := findCaller(hook.Skip)
	entry.Data[hook.FileName] = fileName
	entry.Data[hook.Line] = line
	return nil
}

//实现 logrus.Hook 接口
func (hook *myHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

//自定义hook
func NewMyHook(skip int, levels ...logrus.Level) logrus.Hook {
	hook := myHook{
		FileName: "file",
		Line:     "line",
		Skip:     skip,
		levels:   levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if string(file[i]) == "/" {
			n++
			if n >= 2 {
				//fmt.Println(n >= 2, file)
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}

func findCaller(skip int) (string, int) {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") && !strings.HasPrefix(file, "logger") {
			break
		}
	}
	return file, line
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "20060102 15:04:05"})
}

func EnableDebug() {
	logrus.SetLevel(logrus.DebugLevel)
}

func EnableFileLine(skip int) {
	logrus.AddHook(NewMyHook(skip))
}

// rotate log to file
func LogToFile(file string) {
	path := file
	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)
	if err != nil {
		logrus.Errorf("can't create rotatelogs")
		return
	}

	logrus.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: writer,
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
		},
		&logrus.JSONFormatter{},
	))

}

// Error .
func Error(v ...interface{}) {
	logrus.Error(v...)
}

// Errorf .
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Fatal .
func Fatal(v ...interface{}) {
	logrus.Panic(v...)
}

// Fatalf .
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

// Info .
func Info(v ...interface{}) {
	logrus.Info(v...)
}

// Infof .
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Debug .
func Debug(v ...interface{}) {
	logrus.Debug(v...)
}

// Debugf .
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}
