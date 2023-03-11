package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"

	"log"
	"os"
)

func errCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loadJpeg(filename string, ext []string) (image.Image, error) {

	f, err := os.Open(filename) //открываем кадр
	errCheck(err)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	//Декодируем как jpg
	if ext[1] == "jpg" || ext[1] == "jpeg" {
		img, err := jpeg.Decode(f)
		errCheck(err)
		return img, nil
	}

	//Декодируем как png
	if ext[1] == "png" {
		img, err := png.Decode(f)
		errCheck(err)
		return img, nil
	}

	return nil, nil
}
func diff(a, b uint32) int64 { //Сравнение количества пикселей кадра определенного цвета
	if a > b {
		return int64(a - b)
	}
	return int64(b - a)
}
func extCheck(ext []string){
	if ext[1] != "jpg" || ext[1] != "jpeg" || ext[1] != "png" {
		log.Printf("Формат " + ext[1] + " не поддерживается")
		os.Exit(1)
	}
}
func main() {
	if len(os.Args) != 3 {
		log.Printf("Необходимо использовать 2 параметра <Полное_название_дирректории_источника> <Полный_путь_до_дирректории_для_результата>\"")
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

	j := 1
	dirName := "№" + strconv.Itoa(j)
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

		firstPic, err := loadJpeg(name1, ext1) //Распаковываем первый кадр
		errCheck(err)

		secondPic, err := loadJpeg(name2, ext2) //Распаковываем второй кадр
		errCheck(err)

		if firstPic.ColorModel() != secondPic.ColorModel() { //Проверка на цветовую модель
			log.Fatal("different color models")
		}

		b := firstPic.Bounds() //Проверка на одинаковость размера кадров
		if !b.Eq(secondPic.Bounds()) {
			log.Fatal("different image sizes")
		}

		var sum int64
		for y := b.Min.Y; y < b.Max.Y; y++ { //Цикл от минимального до максимального Y
			for x := b.Min.X; x < b.Max.X; x++ { //Цикл от минимального до максимального X
				r1, g1, b1, _ := firstPic.At(x, y).RGBA()
				r2, g2, b2, _ := secondPic.At(x, y).RGBA()
				sum += diff(r1, r2) //Смотрим разницу красного у первого и второго кадра
				sum += diff(g1, g2) //Смотрим разницу зеленого у первого и второго кадра
				sum += diff(b1, b2) //Смотрим разницу синего у первого и второго кадра
				//Альфа не учитывается поскольку указывает на непрозрачность пикселя от 0 до 1
			}
		}

		amountOfPixels := (b.Max.X - b.Min.X) * (b.Max.Y - b.Min.Y)         //Сколько всего пикселе (по сути площадь кадра)
		result := float64(sum*100) / (float64(amountOfPixels) * 0xffff * 3) //Превращаем разницу в проценты
		fmt.Printf("%d and %d Image difference: %f%%\n", i+1, i+2, result)

		if result <= 0.5 { // Если коэф подобия меньше 0.5%

			oldLocation := "./" + name1                             //Путь к кадру
			newLocation := outputPath + "/" + dirName + "/" + name1 //Новый путь к кадру
			err := os.Rename(oldLocation, newLocation)              // Переносим кадр
			errCheck(err)
			if i+2 == amountOfFiles { // При сравнении предпоследнего и последнего кадра получили коэф подобия меньше 0.5%
				//Копируем второй кад вместе в первым
				oldLocation = "./" + name2
				newLocation = outputPath + "/" + dirName + "/" + name2
				err := os.Rename(oldLocation, newLocation)
				errCheck(err)
				break
			}

		}
		if result > 0.5 { // Если коэф подобия больше 0.5%
			oldLocation := "./" + name1                             //Путь к кадру
			newLocation := outputPath + "/" + dirName + "/" + name1 //Новый путь к кадру
			err := os.Rename(oldLocation, newLocation)              // Переносим кадр
			errCheck(err)

			j++
			dirName = "№" + strconv.Itoa(j) //Поскольку второй кадру отличаектся от первого, то надо создать для него новую диррекеторию
			err1 := os.Mkdir(outputPath+"/"+dirName, 0750)
			if err1 != nil && !os.IsExist(err1) {
				log.Fatal(err1)
			}
			if err1 != nil && os.IsExist(err1) {
				log.Printf("File " + dirName + " is already exist")
			}
			if i+2 == amountOfFiles {
				// При сравнении предпоследнего и последнего кадра получили коэф подобия больше 0.5%
				//Копируем второй кад вместе в первым
				oldLocation = "./" + name2
				newLocation = outputPath + "/" + dirName + "/" + name2
				err := os.Rename(oldLocation, newLocation)
				errCheck(err)
				break
			}

		}
	}

}
