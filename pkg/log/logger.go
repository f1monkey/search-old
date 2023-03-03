package log

import (
	"errors"
	"os"

	"github.com/f1monkey/search/pkg/errs"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// New create new logger instance
func New(level zapcore.Level, traceLevel zapcore.Level, name string, version string, instanceID string, env string) (*zap.Logger, error) {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		NameKey:        "app",
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "caller",
		FunctionKey:    "function",
		MessageKey:     "msg",
		StacktraceKey:  "trace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	return zap.New(
		zapcore.NewCore(&stackTraceEncoder{encoder}, os.Stdout, level),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Fields(
			zap.String("version", version),
			zap.String("app", name),
			zap.String("instance", instanceID),
			zap.String("env", env),
		),
	), nil
}

type stackTraceEncoder struct {
	zapcore.Encoder
}

func (w *stackTraceEncoder) Clone() zapcore.Encoder {
	return &stackTraceEncoder{w.Encoder.Clone()}
}

func (w *stackTraceEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	if len(fields) == 0 {
		return w.Encoder.EncodeEntry(entry, fields)
	}
	for i := range fields {
		if fields[i].Type == zapcore.ErrorType {
			err, ok := fields[i].Interface.(error)
			if !ok {
				continue
			}
			var traceableErr *errs.Error
			if errors.As(err, &traceableErr) {
				fields = append(fields, zap.String("trace", string(traceableErr.StackTrace())))
			}
		}
	}
	return w.Encoder.EncodeEntry(entry, fields)
}

func NewNop() *zap.Logger {
	return zap.NewNop()
}
