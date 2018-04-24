package exhibition

import (
	"database/sql"
	"time"
)

type Exhibition struct {
	ID         string
	Pathname   string
	Hash       string
	CommitTime time.Time
	CORS       sql.NullString
}

func (e *Exhibition) GetBlobPath() string {
	return e.ID + "/" + e.Hash
}
