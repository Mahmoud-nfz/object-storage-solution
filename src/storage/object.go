package storage

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"
)

func MakeObject(bucketName, objectName string, data []byte) error {
	reader := bytes.NewReader(data)
	opts := minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	}
	_, err := MinioClient.PutObject(context.Background(), bucketName, objectName, reader, int64(len(data)), opts)
	if err != nil {
		return err
	}

	return nil
}

func DeleteObject(bucketName, objectName string) error {
	opts := minio.RemoveObjectOptions{}
	err := MinioClient.RemoveObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return err
	}

	return nil
}

func RenameObject(bucketName, oldName, newName string) error {
	src := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: oldName,
	}
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: newName,
	}
	_, err := MinioClient.CopyObject(context.Background(), dst, src)
	if err != nil {
		return err
	}

	opts := minio.RemoveObjectOptions{}
	err = MinioClient.RemoveObject(context.Background(), bucketName, oldName, opts)
	if err != nil {
		return err
	}

	return nil
}

func CopyObjectToBucket(bucketName, destination, objectName string) error {
	src := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: objectName,
	}
	dst := minio.CopyDestOptions{
		Bucket: destination,
		Object: objectName,
	}

	_, err := MinioClient.CopyObject(context.Background(), dst, src)
	if err != nil {
		return err
	}

	return nil
}

