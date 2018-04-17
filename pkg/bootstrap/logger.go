package bootstrap

import (
	"os"

	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
)

var Logger logrus.FieldLogger

func setupLogger() logrus.FieldLogger {
	logger := logrus.New()
	if os.Getenv("VERBOSE") != "" {
		logger.SetLevel(logrus.DebugLevel)
	}
	Logger = logger
	return logger
}

func injectLogger(graph *inject.Graph) {
	logger := setupLogger()
	graph.Provide(
		&inject.Object{Name: "handler logger", Value: logger.WithField("source", "handler")},
		&inject.Object{Name: "exhibition logger", Value: logger.WithField("source", "exhibition")},
		&inject.Object{Name: "artwork logger", Value: logger.WithField("source", "artwork")},
	)
}
