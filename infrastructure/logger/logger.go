package logger

import (
	"log"
	"os"
)

type Logger struct {
}

type LoggerInterface interface {
	logError(statusCode int, err error)
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) LogError(statusCode int, err error) {
	errorLog := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	errorLog.Printf("status code: %d, error: %s", statusCode, err.Error())
}
