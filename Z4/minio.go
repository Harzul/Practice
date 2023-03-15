package main

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

func MinioConnect() (*minio.Client, context.Context) {
	ctx := context.Background()
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	errCheck(err)
	return minioClient, ctx
}

func DownloadVideo(minioClient *minio.Client, ctx context.Context, bucketName string, folderName string) error {
	err := minioClient.FGetObject(ctx, bucketName, folderName+"/"+objectName, "./temp/temp.mp4", minio.GetObjectOptions{})
	if err != nil {
		return errors.New("can't get file, no such directory")
	}
	log.Printf("Successfully scanned %s\n", objectName)

	err = os.Mkdir("./temp", 0750) //Создаем папку с результатом если она не существует
	if err != nil && !os.IsExist(err) {
		return errors.New("can't create directory")
	}
	if err != nil && os.IsExist(err) {
		log.Printf("File " + "temp" + " is already exist")
	}
	err = os.Chdir("temp") //Переходим в эту папку (меняем рабочую дирректорию)
	if err != nil {
		return errors.New("can't enter dirrectory")
	}
	return nil
}
func CreateNewBucket(minioClient *minio.Client, ctx context.Context, bucketName string) error {
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s is already exist \n", bucketName)
		} else {
			return errors.New("bucket creation error")
		}
	}
	return nil
}
