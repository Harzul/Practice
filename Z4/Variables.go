package main

const (
	InputQueueName  = "Input"
	url             = "amqp://rmuser:rmpassword@localhost:5672/"
	endpoint        = "127.0.0.1:9000"
	accessKeyID     = "admin"
	secretAccessKey = "adminadmin"
	objectName      = "VID.mp4"
)

var input = map[string]string{
	"minio_bucket": "",
	"folder_name":  "",
}
var output = map[string]string{
	"minio_bucket": "",
	"folder_name":  "",
}
