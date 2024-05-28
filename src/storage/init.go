package storage

import (
	"data-storage/src/config"

	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func init() {
	useSSL := false
	client, err := minio.New(config.Env.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Env.MinioAccessKey, config.Env.MinioSecretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln("Error initializing Minio client:", err)
	} else {
		MinioClient = client
	}
}
