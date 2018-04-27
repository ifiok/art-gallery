package bootstrap

import (
	"os"

	"golang.ysitd.cloud/blob/cache"
	"golang.ysitd.cloud/blob/minio"

	"github.com/facebookgo/inject"
	"github.com/golang/groupcache"
)

const cacheGroup = "gallery"
const cacheSize = 16 << 20

func setupGroup() *groupcache.Group {
	if group := groupcache.GetGroup(cacheGroup); group != nil {
		return group
	}

	getter := cache.NewGetter(minio.New(setupMinio()), os.Getenv("S3_BUCKET"), Logger.WithField("source", "artwork"))
	return groupcache.NewGroup(cacheGroup, cacheSize, getter)
}

func injectArtworkStore(graph *inject.Graph) {
	graph.Provide(
		&inject.Object{Value: setupGroup()},
	)
}
