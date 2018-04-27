package bootstrap

import (
	"os"

	_ "github.com/lib/pq"

	"golang.ysitd.cloud/db"

	"github.com/facebookgo/inject"
)

func createDB() db.Opener {
	return db.NewOpener("postgres", os.Getenv("DB_URL"))
}

func injectDB(graph *inject.Graph) {
	pool := createDB()
	graph.Provide(&inject.Object{
		Name:  "db",
		Value: pool,
	})
}
