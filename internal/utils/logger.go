package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	// Create a new logger instance
	Log = logrus.New()

	// Set the output to stdout (you can change it to a file or any other output)
	Log.Out = os.Stdout

	// Set the log level (Debug, Info, Warn, Error, Fatal, Panic)
	Log.SetLevel(logrus.DebugLevel)

	// Set a formatter (you can use JSON or text based on your preference)
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Alternatively, you can use the text formatter for more human-readable logs
	// Log.SetFormatter(&logrus.TextFormatter{
	// 	FullTimestamp: true,
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// })
}

func LogInfo(message string, fields logrus.Fields) {
	Log.WithFields(fields).Info(message)
}

func LogWarn(message string, fields logrus.Fields) {
	Log.WithFields(fields).Warn(message)
}

func LogError(err error, message string, fields logrus.Fields) {
	Log.WithFields(fields).WithError(err).Error(message)
}

func LogDebug(message string, fields logrus.Fields) {
	Log.WithFields(fields).Debug(message)
}

func LogFatal(err error, message string, fields logrus.Fields) {
	Log.WithFields(fields).WithError(err).Fatal(message)
}
