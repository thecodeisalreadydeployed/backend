package workloadcontroller

import "go.uber.org/zap"

func NewDeployment(appID string) error {
	logger := zap.L().Sugar().With("appID", appID)
	_ = logger
	return nil
}
