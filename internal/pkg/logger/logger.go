package logger

import "go.uber.org/zap"

type Logger = zap.Logger
type Field = zap.Field

var (
	String = zap.String
	Int    = zap.Int
	Error  = zap.Error
	Any    = zap.Any
)

func New(env string) *zap.Logger {
	if env == "prod" || env == "production" {
		log, _ := zap.NewProduction()
		return log
	}
	log, _ := zap.NewDevelopment()
	return log
}
