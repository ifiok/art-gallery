package bootstrap

import (
	"os"

	_ "github.com/lib/pq"

	"code.ysitd.cloud/common/go/db"
	"github.com/facebookgo/inject"
)

func createDB() db.Pool {
	return db.NewPool("postgres", os.Getenv("DB_URL"))
}

func injectDB(graph *inject.Graph) {
	pool := createDB()
	graph.Provide(&inject.Object{
		Name:  "db",
		Value: pool,
	})
}
