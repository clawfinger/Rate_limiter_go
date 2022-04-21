package logger

import (
	"fmt"
	"os"

	"github.com/clawfinger/ratelimiter/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
}

type loggerImpl struct { // TODO
	zapLogger *zap.SugaredLogger
}

func New(cfg *config.Config) (Logger, error) {
	zapLevel, err := loggerLevelFromString(cfg.Data.Logger.Level)
	if err != nil {
		return nil, fmt.Errorf("Error on logger init, Reason: %w", err)
	}

	pe := zap.NewProductionEncoderConfig()

	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	fileSync, _, err := zap.Open(cfg.Data.Logger.Filename)
	if err != nil {
		return nil, fmt.Errorf("Error on logger init, Reason: %w", err)
	}
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(fileSync), zapLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLevel),
	)

	l := zap.New(core)

	return &loggerImpl{
		zapLogger: l.Sugar(),
	}, nil
}

func (l *loggerImpl) Info(args ...interface{}) {
	l.zapLogger.Info(args)
}

func (l *loggerImpl) Debug(args ...interface{}) {
	l.zapLogger.Info(args)
}

func (l *loggerImpl) Error(args ...interface{}) {
	l.zapLogger.Info(args)
}

func loggerLevelFromString(lvl string) (zapcore.Level, error) {
	switch lvl {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown logger level %s, info set es default", lvl)
	}
}
