package bootstrap

import (
	"os"

	"code.ysitd.cloud/component/art/gallery/pkg/modals/artwork"
	"code.ysitd.cloud/toolkit/blob/cache"
	"code.ysitd.cloud/toolkit/blob/client/minio"

	"github.com/facebookgo/inject"
	"github.com/golang/groupcache"
)

const cacheGroup = "gallery"
const cacheSize = 16 << 20

func setupBlobStore() *cache.CachedBlobStore {
	if group := groupcache.GetGroup(cacheGroup); group != nil {
		return &cache.CachedBlobStore{
			Group: group,
		}
	}

	getter := cache.NewGetter(&minio.Store{Client: setupMinio()}, os.Getenv("S3_BUCKET"), Logger.WithField("source", "artwork"))
	group := groupcache.NewGroup(cacheGroup, cacheSize, getter)
	return &cache.CachedBlobStore{
		Group: group,
	}
}

func setupArtworkStore() *artwork.Store {
	blob := setupBlobStore()
	return &artwork.Store{
		CachedBlobStore: blob,
	}
}

func injectArtworkStore(graph *inject.Graph) {
	graph.Provide(
		&inject.Object{Value: setupArtworkStore()},
	)
}
