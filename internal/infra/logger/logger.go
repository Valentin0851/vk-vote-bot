package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	// ... другие методы
}

func New(level, format string) *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)

	switch format {
	case "text":
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
	default:
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		log.SetLevel(logrus.InfoLevel)
	}
	log.SetLevel(lvl)

	return log
}
