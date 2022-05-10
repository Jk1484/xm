package logger

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(New)

type Logger interface {
	Logger() *zap.SugaredLogger
}

type logger struct {
	logger *zap.SugaredLogger
}

type Params struct {
	fx.In
}

func New(p Params) Logger {
	l, _ := zap.NewProduction()
	return &logger{
		logger: l.Sugar(),
	}
}

func (l *logger) Logger() *zap.SugaredLogger {
	return l.logger
}
