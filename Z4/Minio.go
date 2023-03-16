package main

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

func (m *Minio) Connect() {
	(*m).ctx = context.Background()
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	errCheck(err)
	(*m).client = client
}

func (m *Minio) DownloadVideo(bucketName string, folderName string) error {
	err := (*m).client.FGetObject((*m).ctx, bucketName, folderName+"/"+objectName, "./temp/temp.mp4", minio.GetObjectOptions{})
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
func (m *Minio) CreateNewBucket(bucketName string) error {
	err := (*m).client.MakeBucket((*m).ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := (*m).client.BucketExists((*m).ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s is already exist \n", bucketName)
		} else {
			return errors.New("bucket creation error")
		}
	}
	return nil
}
