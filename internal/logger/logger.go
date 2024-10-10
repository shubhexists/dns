package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitializeLogger() {
	Log = logrus.New()

	Log.Out = os.Stdout

	Log.SetLevel(logrus.DebugLevel)

	Log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		PadLevelText:    true,
	})
}
