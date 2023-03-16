package main

const (
	InputQueueName  = "Input" //название очереди откуда получаем задачи
	url             = "amqp://rmuser:rmpassword@localhost:5672/"
	endpoint        = "127.0.0.1:9000"
	accessKeyID     = "admin"      //minio login
	secretAccessKey = "adminadmin" //minio password
	objectName      = "VID.mp4"    //название видеофайла
	goalBucket      = "frames"     //название ведра с результатом
	folder          = "sameframes" // Название папки с кадрами
)

var input = map[string]string{
	"minio_bucket": "",
	"folder_name":  "",
}
var output = map[string]string{
	"minio_bucket": "",
	"folder_name":  "",
}
