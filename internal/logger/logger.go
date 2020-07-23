package logger

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	kitzap "github.com/go-kit/kit/log/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Create to be used.
func Create(level string) log.Logger {
	if strings.ToLower(level) == "nil" {
		return Nil()
	}

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	logger := zap.New(zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), zapLevel(level)))

	return kitzap.NewZapSugarLogger(logger, zapLevel(level))
}

// Nil returns a nil logger
func Nil() log.Logger {
	return log.NewLogfmtLogger(ioutil.Discard)
}

func zapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "err", "error":
		return zapcore.ErrorLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
