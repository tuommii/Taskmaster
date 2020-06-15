package logger

import (
	"io"
	"log"
	"log/syslog"
	"os"
	"sync"
)

var (
	once   sync.Once
	logger *syslog.Writer
	warnL  *log.Logger
	infoL  *log.Logger
	errL   *log.Logger
)

// ChangeOutput ...
func ChangeOutput(file io.Writer) {
	infoL = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	warnL = log.New(file, "WARNING: ", log.Ldate|log.Ltime)
	errL = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
}

func init() {
	file, err := os.OpenFile("/tmp/taskmaster.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	infoL = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	warnL = log.New(file, "WARNING: ", log.Ldate|log.Ltime)
	errL = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
}

// Info ..
func Info(elem ...interface{}) {
	infoL.Println(elem...)
}

// Warning ..
func Warning(elem ...interface{}) {
	warnL.Println(elem...)
}

// Error ..
func Error(elem ...interface{}) {
	errL.Println(elem...)
}

// Fatal ..
func Fatal(elem ...interface{}) {
	errL.Fatal(elem...)
}

// Get returns logger instance, singleton
func Get(useSyslog bool) *syslog.Writer {
	var err error
	once.Do(func() {
		logger, err = syslog.New(syslog.LOG_NOTICE, "taskmaster")
		if err == nil {
			log.SetOutput(logger)
		}
	})
	return logger
}
