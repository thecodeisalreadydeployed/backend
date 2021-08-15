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

func Info(message string, fields ...zap.Field) {
	instance.Info(message, fields...)
}
