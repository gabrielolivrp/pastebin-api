package logging

import (
	"errors"

	"github.com/gabrielolivrp/pastebin-api/pkg/config"
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value interface{}
}

type logger struct {
	logger *zap.Logger
}

func NewLogger(env config.Environment) (Logger, error) {
	var l *zap.Logger
	var err error
	switch env {
	case config.Development:
		l, err = zap.NewDevelopment()
	case config.Production:
		l, err = zap.NewProduction()
	case config.Test:
		l, err = zap.NewNop(), nil
	}

	if l == nil {
		return nil, errors.New("failed to create logger: logger is nil")
	}

	if err != nil {
		return nil, err
	}
	return &logger{logger: l}, nil
}

func (l *logger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, toZapFields(fields)...)
}

func (l *logger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, toZapFields(fields)...)
}

func (l *logger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, toZapFields(fields)...)
}

func (l *logger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, toZapFields(fields)...)
}

func (l *logger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, toZapFields(fields)...)
}

func (l *logger) With(fields ...Field) *logger {
	return &logger{l.logger.With(toZapFields(fields)...)}
}

func toZapFields(fields []Field) []zap.Field {
	var zapFields []zap.Field
	for i := 0; i < len(fields); i += 1 {
		zapFields = append(zapFields, zap.Any(fields[i].Key, fields[i].Value))
	}
	return zapFields
}
