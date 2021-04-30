package logger

import (
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogConfig() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.MessageKey = "message"

	config := zap.NewProductionConfig()
	var logLevel zapcore.Level
	switch viper.GetString("LOG.LEVEL") {
	case "info":
		logLevel = zapcore.InfoLevel
	case "debug":
		logLevel = zapcore.DebugLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		log.Fatal("There is no log level config")
	}
	config.Level = zap.NewAtomicLevelAt(logLevel)
	if viper.GetString("LOG.ENV") == "dev" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Encoding = "console"
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		config.Encoding = "json"
	}
	config.EncoderConfig = encoderConfig

	logger, _ := config.Build()
	return logger
}
