package main

import (
	"net/http"
	"os"
	"time"

	"code.ysitd.cloud/component/art/gallery/pkg/bootstrap"
	"code.ysitd.cloud/component/art/gallery/pkg/bootstrap/groupcache"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
)

const groupUpdateInterval = time.Minute * 5

func main() {
	logger := bootstrap.Logger
	mainLogger := logger.WithField("source", "main")
	if !groupcache.SingleNode {
		go groupcache.Listen(logger.WithField("source", "groupcache_http"))
		go func(logger logrus.FieldLogger) {
			time.Sleep(groupUpdateInterval)
			groupcache.UpdatePeer(logger)
		}(logger.WithField("source", "groupcache_update"))
	}

	handler := bootstrap.GetHandler()

	httpLog := logger.WithField("source", "http")
	handler = handlers.CombinedLoggingHandler(httpLog.Writer(), handler)

	recoverLogger := logger.WithField("source", "recover")
	handler = handlers.RecoveryHandler(
		handlers.RecoveryLogger(recoverLogger),
		handlers.PrintRecoveryStack(true),
	)(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mainLogger.Infof("Listen at %s", port)

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		mainLogger.Fatalln(err)
	}
}
