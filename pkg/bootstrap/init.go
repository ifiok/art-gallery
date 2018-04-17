package bootstrap

import (
	"net/http"

	"code.ysitd.cloud/component/art/gallery/pkg/service"
	"github.com/facebookgo/inject"
)

var handler http.Handler

func init() {
	var graph inject.Graph

	injectLogger(&graph)
	injectCache(&graph)
	injectDB(&graph)
	injectMinio(&graph)

	handler = new(service.Handler)
	graph.Provide(&inject.Object{Value: handler})
}

func GetHandler() http.Handler {
	return handler
}
