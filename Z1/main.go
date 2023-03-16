package main

import (
	"fmt"
	"io"
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
func extCheck(ext string) {
	if ext != "jpg" && ext != "jpeg" && ext != "png" {
		log.Printf("Формат " + ext + " не поддерживается")
		os.Exit(1)
	}
}
func ScanFile(fileName string) []byte {
	original, err := os.Open(fileName) //Открываем видео
	errCheck(err)
	defer func(original *os.File) {
		err := original.Close()
		errCheck(err)
	}(original)

	data := make([]byte, 0) //Куда читаем данные
	tData := make([]byte, 1024)
	for { //Цикл считывания байт
		n, err := original.Read(tData)
		for index, _ := range tData {
			data = append(data, tData[index])
		}

		if err == io.EOF { //Если конец файла = конец
			break
		}
		_ = n
	}
	return data
}
func MakeDir(dirName string) {
	err := os.Mkdir(dirName, 0750) //Создаем папку с результатом если она не существует
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	if err != nil && os.IsExist(err) {
		log.Printf("File " + dirName + " is already exist")
	}

	err = os.Chdir(dirName) //Переходим в эту папку (меняем рабочую дирректорию)
	errCheck(err)
}
func CreateFile(fileName string, data []byte) {
	tempFile, err := os.Create(fileName) //Создаем файл с названием видеофайла
	errCheck(err)
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		errCheck(err)
	}(tempFile)

	write, err := tempFile.Write(data)
	errCheck(err)
	_ = write
}
func RunFfmpeg(fileName string, ext string) {
	s := "ffmpeg -i " + fileName + " image%08d." + ext
	args := strings.Split(s, " ")

	cmd := exec.Command(args[0], args[1:]...) //Запускаем ffmpeg
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Running ffmpeg failed: %v", err)
	}
}
func Clear(fileName string) {
	err := os.Remove(fileName) //удаляем копию видео
	errCheck(err)
}
func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Необходимо использовать 3 параметра <Полное_название_файла_с_расширением> <Полный_путь_до_дирректории_куда_класть_кадры> <Формат_сохранения_кадров>")
		return
	}
	fileName := os.Args[1] //Название видеофайла
	dirName := os.Args[2]  //Полный путь к дирректории
	ext := os.Args[3]      //Расширение кадров

	extCheck(ext)
	data := ScanFile(fileName)
	MakeDir(dirName)

	CreateFile(fileName, data)
	RunFfmpeg(fileName, ext)

	Clear(fileName)
}

