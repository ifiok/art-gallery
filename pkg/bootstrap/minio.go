package bootstrap

import (
	"os"

	"github.com/facebookgo/inject"
	"github.com/minio/minio-go"
)

func setupMinio() (client *minio.Client) {
	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		endpoint = "s3.amazoneaws.com"
	}
	client, _ = minio.New(
		endpoint,
		os.Getenv("S3_ACCESS_KEY_ID"),
		os.Getenv("S3_SECRET_ACCESS_KEY"),
		os.Getenv("S3_INSCRUE") == "",
	)
	return
}

func injectMinio(graph *inject.Graph) {
	bucket := os.Getenv("S3_BUCKET")
	graph.Provide(
		&inject.Object{Value: setupMinio()},
		&inject.Object{Name: "bucket", Value: &bucket},
	)
}
