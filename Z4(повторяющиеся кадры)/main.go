package main

import (
	"log"
	"strconv"
)

func main() {
	minioClient, minCtx := MinioConnect()
	ch, conn := RabbitConnect()
	defer func() {
		_ = ch.Close()
		_ = conn.Close()
	}()

	pwd := GetPwd()
	messages := CreateConsumer(ch)

	var forever chan struct{}
	var counter = 0
	go func() {
		rabCtx := RabCtxDefine()

		for message := range messages {

			str := message.Body
			err := ParseString(str)
			if err != nil {
				log.Println(">>>>>>>>>> ", err)
				continue
			}

			bucketName := input["minio_bucket"]
			folderName := input["folder_name"]

			err = DownloadVideo(minioClient, minCtx, bucketName, folderName)
			if err != nil {
				log.Println(">>>>>>>>>> ", err)
				continue
			}
			RunFfmpeg()
			files, amountOfFiles := GetDirInfo()

			output["minio_bucket"] = goalBucket
			output["folder_name"] = folder + strconv.Itoa(counter)

			err = CreateNewBucket(minioClient, minCtx, output["minio_bucket"])
			if err != nil {
				log.Println(">>>>>>>>>> ", err)
				continue
			}
			ImageCompare(amountOfFiles, files, minioClient, minCtx, output["minio_bucket"], output["folder_name"])
			body := "minio_bucket : " + output["minio_bucket"] + ",\n" + "folder_name : " + output["folder_name"]

			SendMessage(ch, rabCtx, body)
			clear(pwd)
			counter++
		}

	}()
	<-forever

}
