package logger

import (
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
	if err == nil {
		logFilePath := filepath.Join(appDataDir, consts.AppDataDir, consts.LogsDir, "errors.log")

    lumberjackSink := zapcore.AddSync(&lumberjack.Logger{
			Filename:		logFilePath,
			MaxSize:		100,
			MaxBackups:	3,
			MaxAge:			7,
			Compress:		true,
    })

    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    // Configure core
    core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			lumberjackSink,
			zapcore.ErrorLevel,
    )

		logger := zap.New(core)
		defer logger.Sync()

		Logger = logger.Sugar()

		return
	}

	// Fallback to console logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	Logger = logger.Sugar()
}
