package main

import (
	"net/http"
	"os"
	"time"

	"code.ysitd.cloud/component/art/gallery/pkg/bootstrap"
	"code.ysitd.cloud/component/art/gallery/pkg/bootstrap/groupcache"
	"github.com/sirupsen/logrus"
)

const groupUpdateInterval = time.Minute * 5

func main() {
	if !groupcache.SingleNode {
		go groupcache.Listen(bootstrap.Logger.WithField("source", "groupcache_http"))
		go func(logger logrus.FieldLogger) {
			time.Sleep(groupUpdateInterval)
			groupcache.UpdatePeer(logger)
		}(bootstrap.Logger.WithField("source", "groupcache_update"))
	}

	handler := bootstrap.GetHandler()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, handler)
}
