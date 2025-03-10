package logger

import "go.uber.org/zap"

var Sugar *zap.SugaredLogger

func InitLogger() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	Sugar = logger.Sugar()
}
