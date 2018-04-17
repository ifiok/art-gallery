package exhibition

import "time"

type Exhibition struct {
	ID         string
	Pathname   string
	Hash       string
	CommitTime time.Time
}

func (e *Exhibition) GetBlobPath() string {
	return e.ID + "/" + e.Hash
}
