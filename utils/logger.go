package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	Log.SetOutput(os.Stdout)

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	env := os.Getenv("APP_ENV")
	logLevel := os.Getenv("LOG_LEVEL")

	if logLevel == "" {
		if env == "production" {
			Log.SetLevel(logrus.InfoLevel)
		} else {
			Log.SetLevel(logrus.DebugLevel)
		}
	} else {
		lvl, err := logrus.ParseLevel(logLevel)
		if err != nil {
			Log.SetLevel(logrus.InfoLevel)
		} else {
			Log.SetLevel(lvl)
		}
	}

}
