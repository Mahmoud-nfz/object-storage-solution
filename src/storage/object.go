package storage

import (
	"bytes"
	"context"
	"log"

	"github.com/minio/minio-go/v7"
)

func GetObject(bucketName, objectName string) (*minio.Object, error) {
	object, err := MinioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return object, nil
}

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

func ConcatenateObjects(dst minio.CopyDestOptions, srcs ...minio.CopySrcOptions) (minio.UploadInfo, error) {
	uploadInfo, err := MinioClient.ComposeObject(context.Background(), dst, srcs...)
	if err != nil {
		log.Printf("failed to compose object: %v", err)
		return uploadInfo, err
	}

	return uploadInfo, nil
}
