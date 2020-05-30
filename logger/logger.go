package logger

import (
	"log"
	"log/syslog"
	"sync"
)

var (
	once   sync.Once
	logger *syslog.Writer
)

// Get returns logger instance, singleton
func Get() *syslog.Writer {
	var err error
	once.Do(func() {
		logger, err = syslog.New(syslog.LOG_NOTICE, "taskmaster")
		if err == nil {
			log.SetOutput(logger)
		}
	})
	return logger
}
