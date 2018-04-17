package groupcache

import (
	"net/http"

	"github.com/golang/groupcache"
	"github.com/sirupsen/logrus"
)

func Listen(logger logrus.FieldLogger) {
	pool := picker.(*groupcache.HTTPPool)
	mux := http.NewServeMux()
	mux.Handle("/_groupcache/", pool)
	if err := http.ListenAndServe(":50005", mux); err != nil {
		logger.Errorln(err)
	}
}
