package artwork

import (
	"bytes"
	"context"
	"io"
	"time"

	"code.ysitd.cloud/component/art/gallery/pkg/modals/exhibition"
	"github.com/golang/groupcache"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
)

const minioTimout = 30 * time.Second

const cacheGroup = "gallery"
const cacheSize = 16 << 20

type Store struct {
	Group  *groupcache.Group
	Client *minio.Client      `inject:""`
	Bucket string             `inject:"bucket"`
	Logger logrus.FieldLogger `inject:"artwork logger"`
}

func (s *Store) getGroup() *groupcache.Group {
	if s.Group != nil {
		return s.Group
	}

	s.Group = groupcache.GetGroup(cacheGroup)
	if s.Group != nil {
		return s.Group
	}

	s.Group = groupcache.NewGroup(cacheGroup, cacheSize, s)

	return s.Group
}

func (s *Store) Get(gctx groupcache.Context, key string, dest groupcache.Sink) (err error) {
	var ctx context.Context
	if cast, ok := gctx.(context.Context); ok {
		ctx = cast
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), minioTimout)
		defer cancel()
	}

	obj, err := s.Client.GetObjectWithContext(ctx, s.Bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return
	}

	defer obj.Close()

	var buffer bytes.Buffer

	_, err = io.Copy(&buffer, obj)
	if err != nil {
		return
	}

	return dest.SetBytes(buffer.Bytes())
}

func (s *Store) GetWithExhibition(ctx context.Context, e *exhibition.Exhibition) (dest []byte, err error) {
	group := s.getGroup()
	err = group.Get(ctx, e.GetBlobPath(), groupcache.AllocatingByteSliceSink(&dest))
	return
}
