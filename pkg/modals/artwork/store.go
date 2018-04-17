package artwork

import (
	"context"

	"code.ysitd.cloud/component/art/gallery/pkg/modals/exhibition"
	"code.ysitd.cloud/toolkit/blob/cache"
)

type Store struct {
	*cache.CachedBlobStore
}

func (s *Store) GetWithExhibition(ctx context.Context, e *exhibition.Exhibition) (dest []byte, err error) {
	return s.GetBlobWithContext(ctx, e.GetBlobPath())
}
