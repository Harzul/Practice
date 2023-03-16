package main

import (
	"github.com/vitali-fedulov/images4"
	"log"
	"os"
	"strings"
)

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func extCheck(ext []string) {
	if ext[1] != "jpg" && ext[1] != "jpeg" && ext[1] != "png" {
		log.Printf("Формат " + ext[1] + " не поддерживается")
		os.Exit(1)
	}
}
func ScanDir(inputPath string) []os.FileInfo {
	err := os.Chdir(inputPath) //Переходим в дирректорию c изображениями
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
func MakeDir(outputPath string, dirName string) {
	err := os.Mkdir(outputPath+"/"+dirName, 0750) //Пытаемся создать дирректорию куда будет класться результат
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	if err != nil && os.IsExist(err) {
		log.Printf("File is already exist")
	}
}
func ImageCompare(name1 string, name2 string) float64 {
	firstPic, _ := images4.Open(name1)
	secondPic, _ := images4.Open(name2)

	icon1 := images4.Icon(firstPic)
	icon2 := images4.Icon(secondPic)
	r, g, b := images4.EucMetric(icon1, icon2)
	EucResult := r + g + b
	return EucResult
}
func CopyPic(name string, outputPath string, dirName string) {
	oldLocation := "./" + name                             //Путь к кадру
	newLocation := outputPath + "/" + dirName + "/" + name //Новый путь к кадру
	err := os.Rename(oldLocation, newLocation)             // Переносим кадр
	errCheck(err)
}
func main() {
	if len(os.Args) != 3 {
		log.Printf("Необходимо использовать 2 параметра <Полный_путь_до_дирректории_источника> <Полный_путь_до_дирректории_для_результата>\"")
	}
	inputPath := os.Args[1]  //Полный путь до дирректории источника
	outputPath := os.Args[2] //Полный путь до дирректории с результатом

	files := ScanDir(inputPath)
	amountOfFiles := len(files) //Крличество файлов в папке

	MakeDir(outputPath, "")
	dirName := "uniqueFrames"
	MakeDir(outputPath, dirName)

	for i := 0; i < amountOfFiles; i++ { //Цикл сравнения фото

		name1 := files[i].Name()   //Имя первого кадра
		name2 := files[i+1].Name() //Имя второго кадра

		ext1 := strings.Split(name1, ".") //Расширение первого кадра
		extCheck(ext1)
		ext2 := strings.Split(name2, ".") //Расширение второго кадра
		extCheck(ext2)

		EucResult := ImageCompare(name1, name2)

		if EucResult < 150000 {
			if i+2 == amountOfFiles {
				break
			}
			if i == 0 {
				CopyPic(name1, outputPath, dirName)
			}
		} else {
			CopyPic(name2, outputPath, dirName)
			i++
			if i+2 == amountOfFiles {
				break
			}
		}
	}

}