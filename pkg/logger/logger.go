package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

// OutputLog is used to output logs to a external .log file
var OutputLog = logrus.New()

func init() {
	file, err := os.OpenFile("logs/api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		logrus.Fatal("Failed to open log file: ", err)
	}
	OutputLog.Out = file
}

// Log is used to output the logs to the console in the development mode
var Log = logrus.New()
