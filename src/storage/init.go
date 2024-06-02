package storage

import (
	"data-storage/src/config"

	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client
var UploadsBucket = "uploads"
var DataBucket = "data"

func init() {
	useSSL := false
	client, err := minio.New(config.Env.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Env.MinioAccessKey, config.Env.MinioSecretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln("Error initializing Minio client: ", err)
	} else {
		MinioClient = client
	}

	if err := MakeBucket(UploadsBucket); err != nil {
		log.Fatalln("Error initializing uploads bucket: ", err)
	}

	if err := MakeBucket(DataBucket); err != nil {
		log.Fatalln("Error initializing data bucket: ", err)
	}
}
