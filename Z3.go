package main

import (
	"fmt"
	"log"
	"os"
)

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func copyPics(amountOfFiles int, files []os.FileInfo, outputDirName string, start int, end int, ch chan string) {
	//Циклом переносим фото
	for i := start; i < end; i++ {
		fileName := files[i].Name()

		oldLocation := "./" + fileName
		newLocation := outputDirName + "/" + fileName
		err := os.Rename(oldLocation, newLocation)
		errCheck(err)
	}
	ch <- "Done"
}

func main() {

	if len(os.Args) != 3 {
		log.Printf("Необходимо использовать 2 параметра <Полное_название_дирректории_откуда_копируем> <Полный_путь_до_дирректории_куда_копируем>\"")
	}
	dirName := os.Args[1]       //Дирректория с кадрами
	outputDirName := os.Args[2] //Дирректория куда копируем

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

	amountOfFiles := len(files) //Количество файлов

	if amountOfFiles == 0 {
		fmt.Println("Directory is empty")
		return
	}

	err1 := os.Mkdir(outputDirName, 0750) //Пытаемся создать папку для результата
	if err1 != nil && !os.IsExist(err1) {
		log.Fatal(err1)
	}
	if err1 != nil && os.IsExist(err1) {
		log.Printf("File " + dirName + " is already exist")
	}

	if amountOfFiles > 100 {
		//Открываем 4 канала
		ch1 := make(chan string)
		ch2 := make(chan string)
		ch3 := make(chan string)
		ch4 := make(chan string)
		var N = amountOfFiles / 4
		//Запускаем 4 горутины
		go copyPics(amountOfFiles, files, outputDirName, 0, N, ch1)
		go copyPics(amountOfFiles, files, outputDirName, N, N*2, ch2)
		go copyPics(amountOfFiles, files, outputDirName, N*2, N*3, ch3)
		go copyPics(amountOfFiles, files, outputDirName, N*3, amountOfFiles, ch4)
		//Ждем пока все каналы получат "Done"
		fmt.Println("First thread:", <-ch1)
		fmt.Println("Second thread:", <-ch2)
		fmt.Println("Third thread:", <-ch3)
		fmt.Println("Fourth thread:", <-ch4)
	}
	if amountOfFiles < 100 {
		ch := make(chan string)
		go copyPics(amountOfFiles, files, outputDirName, 0, amountOfFiles, ch)
		fmt.Println(<-ch)
	}

}
