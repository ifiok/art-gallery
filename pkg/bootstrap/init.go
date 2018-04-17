package bootstrap

import (
	"net/http"

	"code.ysitd.cloud/component/art/gallery/pkg/service"

	"github.com/facebookgo/inject"
)

var handler http.Handler

type graphInjectFunc func(*inject.Graph)

func init() {
	logger := setupLogger()
	graph := inject.Graph{
		Logger: logger,
	}

	h := new(service.Handler)
	graph.Provide(&inject.Object{Value: h})

	injects := []graphInjectFunc{
		injectLogger,
		injectCache,
		injectDB,
		injectArtworkStore,
	}

	for _, injectFn := range injects {
		injectFn(&graph)
	}

	if err := graph.Populate(); err != nil {
		logger.Fatalln(err)
	}

	handler = h
}

func GetHandler() http.Handler {
	return handler
}
