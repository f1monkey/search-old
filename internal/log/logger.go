package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New create logger instance
func New(level string) (*zap.Logger, error) {
	var logLevel zap.AtomicLevel
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encoderConfig.TimeKey = "time"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		os.Stdout,
		logLevel,
	)
	z := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	return z, nil
}
