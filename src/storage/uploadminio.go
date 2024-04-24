package storage

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
)

func UploadToMinioFolder(filePath string, fileName string, bucketName string) error {

	MinioClient, err := InitializeMinioClient()
	if err != nil {
		// Handle the error, perhaps log it or return it
		log.Printf("Error initializing Minio Client: %v", err)
		return err
	}

	if _, err := MinioClient.StatObject(context.Background(), bucketName, fileName, minio.StatObjectOptions{}); err == nil {
		fileName = fileName + "_"
	}

	err = MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := MinioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			return err
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = MinioClient.PutObject(context.Background(), bucketName, fileName, file, fileInfo.Size(), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	log.Println("File uploaded to Minio server successfully")

	return nil
}
