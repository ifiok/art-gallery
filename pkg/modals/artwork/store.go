package artwork

import (
	"context"

	"code.ysitd.cloud/art/gallery/pkg/modals/exhibition"
	"golang.ysitd.cloud/blob/cache"
)

type Store struct {
	cache.CachedBlobStore `inject:"inline"`
}

func (s *Store) GetWithExhibition(ctx context.Context, e *exhibition.Exhibition) (dest []byte, err error) {
	return s.GetBlobWithContext(ctx, e.GetBlobPath())
}
