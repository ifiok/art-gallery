package bootstrap

import (
	"time"

	"github.com/facebookgo/inject"
	"github.com/patrickmn/go-cache"
)

const (
	defaultExpire = time.Minute * 5
	cleanupExpire = time.Minute * 30
)

func createCache() (c *cache.Cache) {
	return cache.New(defaultExpire, cleanupExpire)
}

func injectCache(graph *inject.Graph) {
	graph.Provide(&inject.Object{
		Value: createCache(),
	})
}
