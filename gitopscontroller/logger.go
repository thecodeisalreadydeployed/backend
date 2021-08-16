package gitopscontroller

import (
	"github.com/thecodeisalreadydeployed/logger"
	"go.uber.org/zap"
)

func Info(msg string) {
	logger.Logger().Info(msg, zap.String("package", "gitopscontroller"))
}
