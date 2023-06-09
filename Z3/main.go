package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

const (
	endpoint        = "127.0.0.1:9000"
	accessKeyID     = "admin"
	secretAccessKey = "adminadmin"
)

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func copyPics(files []os.FileInfo, ctx context.Context,
	start int, end int, ch chan string, minioClient *minio.Client, bucketName string) {

	for i := start; i < end; i++ {
		objectName := files[i].Name()
		filePath := "./" + files[i].Name()

		info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	}
	ch <- "Done"
}
func ScanDir(dirName string) []os.FileInfo {
	err := os.Chdir(dirName) //Переходим в диреркторию с кадрами
	errCheck(err)

	dir, err := os.Open(".") // Открываем текущую директорию
	errCheck(err)
	defer func(dir *os.File) {
		err := dir.Close()
		errCheck(err)
	}(dir)

	files, err := dir.Readdir(-1) // Получаем список файлов и папок
	errCheck(err)

	return files
}
func minioStart(bucketName string) (*minio.Client, context.Context) {
	ctx := context.Background()
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	errCheck(err)

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s is already exist \n", bucketName)
		} else {
			log.Fatalln(err)
		}
	}
	return minioClient, ctx
}

func main() {
	if len(os.Args) != 3 {
		log.Printf("Необходимо использовать 2 параметра <Полный_путь_к_дирректории_откуда_копируем> <ведро>\"")
	}
	dirName := os.Args[1]
	bucketName := os.Args[2]

	files := ScanDir(dirName)

	amountOfFiles := len(files) //Количество файлов
	if amountOfFiles == 0 {
		fmt.Println("Directory is empty")
		return
	}

	minioClient, ctx := minioStart(bucketName)

	if amountOfFiles >= 250 {
		//Открываем 4 канала
		ch1 := make(chan string)
		ch2 := make(chan string)
		ch3 := make(chan string)
		ch4 := make(chan string)
		var N = amountOfFiles / 4
		//Запускаем 4 горутины
		go copyPics(files, ctx, 0, N, ch1, minioClient, bucketName)
		go copyPics(files, ctx, N, N*2, ch2, minioClient, bucketName)
		go copyPics(files, ctx, N*2, N*3, ch3, minioClient, bucketName)
		go copyPics(files, ctx, N*3, amountOfFiles, ch4, minioClient, bucketName)
		//Ждем пока все каналы получат "Done"
		fmt.Println("First thread:", <-ch1)
		fmt.Println("Second thread:", <-ch2)
		fmt.Println("Third thread:", <-ch3)
		fmt.Println("Fourth thread:", <-ch4)
	}
	if amountOfFiles < 250 {
		ch := make(chan string)
		go copyPics(files, ctx, 0, amountOfFiles, ch, minioClient, bucketName)
		fmt.Println(<-ch)
	}

}
