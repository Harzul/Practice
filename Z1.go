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

//go run main.go <Полное_название_файла_с_расширением> <Полный_путь_до_дирректории_куда_класть_кадры> <Формат_сохранения_кадров>

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Необходимо использовать 3 параметра <Полное_название_файла_с_расширением> <Полный_путь_до_дирректории_куда_класть_кадры> <Формат_сохранения_кадров>")
		return
	}
	fileName := os.Args[1] //Название видеофайла
	dirName := os.Args[2]  //Полный путь к дирректории
	ext := os.Args[3]      //Расширение кадров

	if ext != "jpg" && ext != "jpeg" && ext != "png" {
		log.Printf("Формат " + ext + " не поддерживается")
		os.Exit(1)
	}

	s := "ffmpeg -i " + fileName + " image%08d." + ext
	args := strings.Split(s, " ")

	original, err := os.Open(fileName) //Открываем видео
	errCheck(err)
	defer func(original *os.File) {
		err := original.Close()
		errCheck(err)
	}(original)
	/*
		file, err := os.Stat(fileName) //Данные о файле
			if err != nil {
				errCheck(err)
			}

		fileSize := file.Size()        //Размер файла в байтах
		data := make([]byte, fileSize) //Куда читаем данные

		//Если закидывать сразу размер все может упасть при больших размерах? Поэтому там ниже еще вариант, более похожий на правду
		n, err := original.Read(data)
		if err == io.EOF { // если конец файла = конец
			log.Fatal(err)
		}
		_ = n
	*/
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

	err1 := os.Mkdir(dirName, 0750) //Создаем папку с результатом если она не существует
	if err1 != nil && !os.IsExist(err1) {
		log.Fatal(err1)
	}
	if err1 != nil && os.IsExist(err1) {
		log.Printf("File " + dirName + " is already exist")
	}

	err2 := os.Chdir(dirName) //Переходим в эту папку (меняем рабочую дирректорию)
	errCheck(err2)

	tempFile, err := os.Create(fileName) //Создаем файл с названием видеофайла
	errCheck(err)
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		errCheck(err)
	}(tempFile)

	write, err := tempFile.Write(data)
	errCheck(err)
	_ = write

	cmd := exec.Command(args[0], args[1:]...) //Запускаем ffmpeg
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Running ffmpeg failed: %v", err)
	}

	err4 := os.Remove(fileName) //удаляем копию видео
	errCheck(err4)

	fmt.Printf("%s\n", b) //Характеристики видео

}