package logger

import (
	"log"
	"os"
	"path/filepath"

	"github.com/X3ne/hsrpc/src/consts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.SugaredLogger

func init() {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("failed to get user config directory: %v", err)
	}

	errorLogFilePath := filepath.Join(appDataDir, consts.AppDataDir, consts.LogsDir, "errors.log")
	errorLumberjackSink := zapcore.AddSync(&lumberjack.Logger{
		Filename:   errorLogFilePath,
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	})

	infoLogFilePath := filepath.Join(appDataDir, consts.AppDataDir, consts.LogsDir, "logs.log")
	infoLumberjackSink := zapcore.AddSync(&lumberjack.Logger{
		Filename:   infoLogFilePath,
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		errorLumberjackSink,
		zapcore.ErrorLevel,
	)

	infoCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		infoLumberjackSink,
		zapcore.InfoLevel,
	)

	core := zapcore.NewTee(errorCore, infoCore)

	logger := zap.New(core)
	defer logger.Sync()

	Logger = logger.Sugar()
}
