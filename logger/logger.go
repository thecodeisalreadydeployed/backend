package logger

import (
	"sync"

	"go.uber.org/zap"
)

var instance *zap.Logger
var once sync.Once

func Init() {
	once.Do(func() {
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		instance = logger
	})
}

func Logger() *zap.Logger {
	if instance == nil {
		Init()
	}

	return instance
}

func Debug(message string, fields ...zap.Field) {
	Logger().Debug(message, fields...)
}

func Info(message string, fields ...zap.Field) {
	Logger().Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	Logger().Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	Logger().Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	Logger().Fatal(message, fields...)
}
