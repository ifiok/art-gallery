package exhibition

import (
	"context"
	"database/sql"

	"code.ysitd.cloud/common/go/db"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type Store struct {
	Pool   db.Pool            `inject:"db"`
	Cache  *cache.Cache       `inject:""`
	Logger logrus.FieldLogger `inject:"exhibition logger"`
}

func (s *Store) GetExhibition(ctx context.Context, hostname, path string) (e *Exhibition, err error) {
	cacheKey := hostname + "/" + path
	if val, hit := s.Cache.Get(cacheKey); hit {
		s.Logger.Debugf("Load %s%s from cache", hostname, path)
		return val.(*Exhibition), nil
	}

	s.Logger.Debugf("Load %s%s from database", hostname, path)
	e, err = s.getFromDB(ctx, hostname, path)
	if err != nil {
		return
	}
	s.Cache.SetDefault(cacheKey, e)
	return
}

func (s *Store) getFromDB(ctx context.Context, hostname, path string) (e *Exhibition, err error) {
	conn, err := s.Pool.Acquire()
	if err != nil {
		return
	}
	defer conn.Close()

	query := `
WITH recursive union_revision AS (
  WITH target_exhibition AS (
    SELECT exhibition.id, revision FROM exhibition
      INNER JOIN exhibition_host ON exhibition_host.exhibition = exhibition.id
    WHERE hostname = $1
  ),
    linked_revision AS (
      SELECT r.id, parent, commit_time, e.id AS exhibtion FROM revision r
      INNER JOIN target_exhibition e ON r.id = e.revision
    )
  SELECT id, parent, commit_time, exhibtion FROM linked_revision
  UNION ALL
  SELECT r.id, r.parent, r.commit_time, r.exhibtion FROM linked_revision r
    INNER JOIN linked_revision p on p.parent = r.id
)
SELECT r.commit_time, tree.pathname, tree.hash, r.exhibtion FROM union_revision r
  INNER JOIN tree on tree.revision = r.id
WHERE tree.pathname = $2
ORDER BY r.commit_time DESC
LIMIT 1
`

	row := conn.QueryRowContext(ctx, query, hostname, path)

	var instance Exhibition

	if err := row.Scan(&instance.CommitTime, &instance.Pathname, &instance.Hash, &instance.ID); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	e = &instance

	return
}
