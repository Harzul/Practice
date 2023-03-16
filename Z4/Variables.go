package main

import (
	"context"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	InputQueueName  = "Input" //название очереди откуда получаем задачи
	url             = "amqp://rmuser:rmpassword@localhost:5672/"
	endpoint        = "127.0.0.1:9000"
	accessKeyID     = "admin"        //minio login
	secretAccessKey = "adminadmin"   //minio password
	objectName      = "VID.mp4"      //название видеофайла
	goalBucket      = "frames"       //название ведра с результатом
	folder          = "uniqueframes" // Название папки с кадрами
)

var forever chan struct{}
var counter = 0
var Input Data

type Data struct {
	MinioBucket string `json:"minio_bucket"`
	FolderName  string `json:"folder_name"`
}

type Minio struct {
	client *minio.Client
	ctx    context.Context
}
type Rabbit struct {
	ch   *amqp.Channel
	conn *amqp.Connection
	ctx  context.Context
}
