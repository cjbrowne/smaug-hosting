package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Setup() {

	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		// default to prod-safe logging level
		lvl = logrus.ErrorLevel
	}
	logrus.SetLevel(lvl)
	logrus.Infof("Log level = %s", lvl)
}
