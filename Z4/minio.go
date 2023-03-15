package main

import (
	"context"
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
func DownloadVideo(minioClient *minio.Client, ctx context.Context, bucketName string, folderName string) {
	err := minioClient.FGetObject(ctx, bucketName, objectName, "./temp/temp.mp4", minio.GetObjectOptions{})
	errCheck(err)
	log.Printf("Successfully scanned %s\n", objectName)

	err = os.Mkdir("./temp", 0750) //Создаем папку с результатом если она не существует
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	if err != nil && os.IsExist(err) {
		log.Printf("File " + "temp" + " is already exist")
	}
	err = os.Chdir("temp") //Переходим в эту папку (меняем рабочую дирректорию)
	errCheck(err)
}
func CreateNewBucket(minioClient *minio.Client, ctx context.Context, bucketName string) {
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s is already exist \n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}
}
