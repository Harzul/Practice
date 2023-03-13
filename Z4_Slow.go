package main

//НЕ ИТОГОВАЯ ВЕРСИЯ
import (
	"bytes"
	"context"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/vitali-fedulov/images4"
	"io"
	"log"
	"os"
	"strconv"
)

// Принимаю от Кролика ведро и папку
// Считываю оттуда видео
// Задание номер 1
// ЗАдание номер 2
// Отправляю кролику "Done" и ведро с папкой уникальных фреймов
const (
	endpoint        = "127.0.0.1:9000"
	accessKeyID     = "admin"
	secretAccessKey = "adminadmin"
	tempVideoName   = "temp.mp4"
	criteria        = 50000
)

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func readFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"loglevel": "fatal", "vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}
func PrepareMinio(bucketName string, objectName string) (context.Context, *minio.Client) {
	ctx := context.Background()
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	errCheck(err)

	err = minioClient.FGetObject(ctx, bucketName, objectName, "./temp/temp.mp4", minio.GetObjectOptions{})
	errCheck(err)
	log.Printf("Successfully scanned %s\n", objectName)
	return ctx, minioClient
}

func MakeDir(path string) {
	err := os.Mkdir("./"+path, 0750) //Создаем папку с результатом если она не существует
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	if err != nil && os.IsExist(err) {
		log.Printf("Dir " + path + " is already exist")
	}
	err = os.Chdir(path) //Переходим в эту папку (меняем рабочую дирректорию)
	errCheck(err)
}
func MakeNewBucket(ctx context.Context, minioClient *minio.Client, MinioDirName string) {
	err := minioClient.MakeBucket(ctx, MinioDirName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, MinioDirName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s is already exist \n", MinioDirName)
		} else {
			log.Fatalln(err)
		}
	}
}
func GetImage(index int, VideoName string, PicName string) {
	reader := readFrameAsJpeg(VideoName, index)
	img, err := imaging.Decode(reader)
	errCheck(err)
	err = imaging.Save(img, PicName)
	errCheck(err)
}

func CompareImage(name1 string, name2 string) float64 {
	firstPic, _ := images4.Open(name1)
	secondPic, _ := images4.Open(name2)

	icon1 := images4.Icon(firstPic)
	icon2 := images4.Icon(secondPic)
	r, g, b := images4.EucMetric(icon1, icon2)
	EucResult := r + g + b
	return EucResult
}
func clean(pwd string) {
	err := os.Remove(tempVideoName) //удаляем копию видео
	errCheck(err)

	err = os.Chdir(pwd)
	errCheck(err)

	err = os.RemoveAll("./temp")
	errCheck(err)
}
func deleteImage(name1 string, name2 string) {
	err := os.Remove(name1)
	errCheck(err)
	err = os.Remove(name2)
	errCheck(err)
}
func main() {

	pwd, err := os.Getwd()
	errCheck(err)

	bucketName := "video"
	videoName := "VID.mp4"
	ctx, minioClient := PrepareMinio(bucketName, videoName)
	MakeDir("temp")

	amountOfFiles := 330 //Крличество файлов в папке

	MinioDirName := "uniqueframes"
	MakeNewBucket(ctx, minioClient, MinioDirName)

	for i := 0; i < amountOfFiles; i++ { //Цикл сравнения фото

		index1 := i
		index2 := i + 1

		name1 := "image" + strconv.Itoa(index1) + ".jpg"
		name2 := "image" + strconv.Itoa(index2) + ".jpg"

		GetImage(index1, tempVideoName, name1)
		GetImage(index2, tempVideoName, name2)

		compareValue := CompareImage(name1, name2)

		if compareValue < criteria {
			if i+2 == amountOfFiles {
				fmt.Printf("Images %d and %d are similar. %f", i, i+1, compareValue)
				fmt.Println()
				break
			}
			fmt.Printf("Images %d and %d are similar. %f", i, i+1, compareValue)
			fmt.Println()
			if i == 0 {
				info, err := minioClient.FPutObject(ctx, MinioDirName, name1, "./"+name1, minio.PutObjectOptions{})
				errCheck(err)
				log.Printf("Successfully uploaded %s of size %d\n", name1, info.Size)
			}
		} else {
			info, err := minioClient.FPutObject(ctx, MinioDirName, name2, "./"+name2, minio.PutObjectOptions{})
			errCheck(err)
			log.Printf("Successfully uploaded %s of size %d\n", name2, info.Size)
			fmt.Printf("Images %d and %d are distinct. %f", i, i+1, compareValue)
			fmt.Println()
			i++
			if i+2 == amountOfFiles {
				break
			}
		}
		deleteImage(name1, name2)

	}
	clean(pwd)
}
