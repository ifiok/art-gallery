# Art Gallery (Frontend Resource Serving)

Art Gallery is micro service for serving frontend resource

## Requirement

- Golang 1.10.1
- [dep](https://github.com/golang/dep)

## Config

| Name | Usage |
|------|-------|
| DB_URL | Connection Url of Postgres |
| VERBOSE |Use verbose logging if exists |
| S3_ENDPOINT | Endpoint of S3 |
| S3_ACCESS_KEY_ID | S3 Access Key ID |
| S3_SECRET_ACCESS_KEY | S3 Secret Access Key |
| S3_INSECURE | Use insecure for S3 connection if exists |
| S3_BUCKET | Name of S3 bucket for storing assets |
