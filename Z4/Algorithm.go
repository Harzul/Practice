package main

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/vitali-fedulov/images4"
	"log"
	"os"
	"os/exec"
	"strings"
)

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func GetPwd() string {
	pwd, err := os.Getwd()
	errCheck(err)
	return pwd
}
func ParseString(str []byte) error {
	strSpace := strings.Split(string(str), " ")
	if len(strSpace) != 5 {
		return errors.New("wrong Input1")
	}

	newStrings := strings.Split(string(str[:]), "\n")
	SS := strings.Replace(newStrings[0], ",", "", -1)
	minioStr := strings.Split(SS, " : ")
	rabbitStr := strings.Split(newStrings[1], " : ")

	if minioStr[0] == "minio_bucket" && rabbitStr[0] == "folder_name" {
		input[minioStr[0]] = minioStr[1]
		input[rabbitStr[0]] = rabbitStr[1]

		return nil
	}

	return errors.New("wrong Input2")
}
func RunFfmpeg() {
	s := "ffmpeg -i " + "temp.mp4" + " image%08d.jpg"
	args := strings.Split(s, " ")

	cmd := exec.Command(args[0], args[1:]...) //Запускаем ffmpeg
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Running ffmpeg failed: %v", err)
	}

	err4 := os.Remove("temp.mp4") //удаляем копию видео
	errCheck(err4)
}

func GetDirInfo() ([]os.FileInfo, int) {
	dir, err := os.Open(".") // Открываем текущую директорию
	errCheck(err)
	defer func(dir *os.File) {
		err := dir.Close()
		errCheck(err)
	}(dir)

	files, err := dir.Readdir(-1) // Получаем список файлов и папок
	errCheck(err)

	amountOfFiles := len(files) //Крличество файлов в папке
	return files, amountOfFiles
}
func ImageCompare(amountOfFiles int, files []os.FileInfo, minioClient *minio.Client, ctx context.Context, bucketName string, folderName string) {
	for i := 0; i < amountOfFiles; i++ { //Цикл сравнения фото

		name1 := files[i].Name()   //Имя первого кадра
		name2 := files[i+1].Name() //Имя второго кадра

		firstPic, _ := images4.Open(name1)
		secondPic, _ := images4.Open(name2)

		icon1 := images4.Icon(firstPic)
		icon2 := images4.Icon(secondPic)
		r, g, b := images4.EucMetric(icon1, icon2)
		EucResult := r + g + b

		if EucResult < 50000 {
			if i+2 == amountOfFiles {
				break
			}
			if i == 0 {
				_, err := minioClient.FPutObject(ctx, bucketName, name1, "./"+name1, minio.PutObjectOptions{})
				if err != nil {
					log.Fatalln(err)
				}
			}
		} else {
			_, err := minioClient.FPutObject(ctx, bucketName, name2, "./"+name2, minio.PutObjectOptions{})
			if err != nil {
				log.Fatalln(err)
			}

			i++
			if i+2 == amountOfFiles {
				break
			}
		}
	}
}
func clear(pwd string) {
	_ = os.Chdir(pwd)
	_ = os.RemoveAll("./temp")
}
