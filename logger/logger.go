package logger

import (
	"io"
	"os"

	"github.com/47Cid/Castle/config"
	"github.com/sirupsen/logrus"
)

var ProxyLog *logrus.Logger
var WAFLog *logrus.Logger

func InitProxyLogger() {
	ProxyLog = logrus.New()

	// Open a file for writing logs
	file, err := os.OpenFile(config.GetProxyLogFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		ProxyLog.Info("Failed to log to file, using default stderr")
		ProxyLog.Out = os.Stdout
	} else {
		ProxyLog.Out = io.MultiWriter(file, os.Stdout)
	}

	// You can set the logging level here
	ProxyLog.SetLevel(logrus.InfoLevel)
}

func InitWAFProxy() {
	WAFLog = logrus.New()

	// Open a file for writing logs
	file, err := os.OpenFile(config.GetWAFLogFile(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		WAFLog.Info("Failed to log to file, using default stderr")
		WAFLog.Out = os.Stdout
	} else {
		WAFLog.Out = io.MultiWriter(file, os.Stdout)
	}

	// You can set the logging level here
	WAFLog.SetLevel(logrus.InfoLevel)
}
