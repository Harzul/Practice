package main

import (
	"fmt"
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
func main() {
	if len(os.Args) != 3 {
		log.Printf("Необходимо использовать 2 параметра <Полный_путь_до_дирректории_источника> <Полный_путь_до_дирректории_для_результата>\"")
	}
	inputPath := os.Args[1]  //Полный путь до дирректории источника
	outputPath := os.Args[2] //Полный путь до дирректории с результатом

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

	amountOfFiles := len(files) //Крличество файлов в папке

	err1 := os.Mkdir(outputPath, 0750) //Пытаемся создать дирректорию куда будет класться результат
	if err1 != nil && !os.IsExist(err1) {
		log.Fatal(err1)
	}
	if err1 != nil && os.IsExist(err1) {
		log.Printf("File is already exist")
	}
	dirName := "uniqueFrames"
	err2 := os.Mkdir(outputPath+"/"+dirName, 0750) //Создаем дирректорию куда будут складываться первые похожие кадры
	if err2 != nil && !os.IsExist(err2) {
		log.Fatal(err2)
	}
	if err2 != nil && os.IsExist(err2) {
		log.Printf("File " + dirName + " is already exist")
	}

	for i := 0; i < amountOfFiles; i++ { //Цикл сравнения фото

		name1 := files[i].Name()   //Имя первого кадра
		name2 := files[i+1].Name() //Имя второго кадра

		ext1 := strings.Split(name1, ".") //Расширение первого кадра
		extCheck(ext1)
		ext2 := strings.Split(name2, ".") //Расширение второго кадра
		extCheck(ext2)

		firstPic, _ := images4.Open(name1)
		secondPic, _ := images4.Open(name2)

		icon1 := images4.Icon(firstPic)
		icon2 := images4.Icon(secondPic)
		r, g, b := images4.EucMetric(icon1, icon2)
		EucResult := r + g + b

		if EucResult < 150000 {
			if i+2 == amountOfFiles {
				fmt.Printf("Images %d and %d are similar. %f", i, i+1, EucResult)
				fmt.Println()
				break
			}
			fmt.Printf("Images %d and %d are similar. %f", i, i+1, EucResult)
			fmt.Println()
			if i == 0 {
				oldLocation := "./" + name1                             //Путь к кадру
				newLocation := outputPath + "/" + dirName + "/" + name1 //Новый путь к кадру
				err := os.Rename(oldLocation, newLocation)              // Переносим кадр
				errCheck(err)
			}
		} else {
			oldLocation := "./" + name2                             //Путь к кадру
			newLocation := outputPath + "/" + dirName + "/" + name2 //Новый путь к кадру
			err := os.Rename(oldLocation, newLocation)              // Переносим кадр
			errCheck(err)
			fmt.Printf("Images %d and %d are distinct. %f", i, i+1, EucResult)
			fmt.Println()
			i++
			if i+2 == amountOfFiles {
				break
			}
		}
	}

}

