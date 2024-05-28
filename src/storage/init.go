package storage

import (
	"data-storage/src/config"

	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitializeMinioClient() (*minio.Client, error) {
	useSSL := false
	MinioClient, err := minio.New(config.Env.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Env.MinioAccessKey, config.Env.MinioSecretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Println("Error initializing Minio client:", err)
		return nil, err
	}

	return MinioClient, nil
}
