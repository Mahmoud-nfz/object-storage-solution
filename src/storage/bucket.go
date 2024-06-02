package storage

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func MakeBucket(destination string) error {
	if found, err := MinioClient.BucketExists(context.Background(), destination); err != nil || !found {
		if errResp := minio.ToErrorResponse(err); !found || errResp.Code == "NoSuchBucket" {
			err = MinioClient.MakeBucket(context.Background(), destination, minio.MakeBucketOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ListBucketObjects(bucketName, prefix string) ([]string, error) {
	objectsChannel := MinioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Recursive: true,
		Prefix: prefix,
	})

	var objects []string
	for object := range objectsChannel {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object.Key)
	}

	return objects, nil
}
