package service

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func PostObject(bucketName string, path string, file multipart.File, fileName string) error {
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(file)
	if err != nil {
		return err
	}

	key := path + fileName

	// MinIOにファイルをアップロード
	_, err = minioClient.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buffer.Bytes()),
		ACL:    aws.String("public-read"), // 必要に応じてACLを設定
	})
	if err != nil {
		return err
	}

	return nil
}

func GetObjects(bucketName string, prefix string) (*s3.ListObjectsV2Output, error) {
	result, err := minioClient.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		fmt.Println("Unable to list buckets", err)
		return nil, err
	}
	return result, nil
}

func GetObjectUrl(bucketName string, key string) (string, error) {
	req, _ := minioClient.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	// 署名付きURLの有効期限を設定
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", err
	}

	return urlStr, nil
}
