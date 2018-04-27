package main

import (
	"net/http"
	"os"
	"time"

	"code.ysitd.cloud/component/art/gallery/pkg/bootstrap"
	"golang.ysitd.cloud/cache/groupcache"
	"golang.ysitd.cloud/cache/groupcache/k8s"

	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
)

const groupUpdateInterval = time.Minute * 5

func init() {
	groupcache.SetSingleNode(os.Getenv("SINGLE_NODE") != "")
	k8s.Setup(2)
}

func main() {
	logger := bootstrap.Logger
	mainLogger := logger.WithField("source", "main")

	if !groupcache.IsSingleNode() {
		go func() {
			if err := groupcache.Listen(); err != nil {
				logger.WithField("source", "groupcache_http").Errorln(err)
			}
		}()
		go func(logger logrus.FieldLogger) {
			for {
				time.Sleep(groupUpdateInterval)
				if err := k8s.Update(); err != nil {
					logger.Errorln(err)
				}
			}
		}(logger.WithField("source", "groupcache_update"))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mainLogger.Infof("Listen at %s", port)

	if err := http.ListenAndServe(":"+port, makeHandler(logger)); err != nil {
		mainLogger.Fatalln(err)
	}
}

func makeHandler(logger logrus.FieldLogger) (handler http.Handler) {
	handler = bootstrap.GetHandler()

	httpLog := logger.WithField("source", "http")
	handler = handlers.CombinedLoggingHandler(httpLog.Writer(), handler)

	recoverLogger := logger.WithField("source", "recover")
	handler = handlers.RecoveryHandler(
		handlers.RecoveryLogger(recoverLogger),
		handlers.PrintRecoveryStack(true),
	)(handler)
	return
}
