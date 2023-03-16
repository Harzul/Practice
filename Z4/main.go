package main

import (
	"encoding/json"
	"log"
	"strconv"
)

func main() {
	m := Minio{}
	m.Connect()
	r := Rabbit{}
	r.RabbitConnect()

	/*defer func() {
		_ = ch.Close()
		_ = conn.Close()
	}()
	*/
	pwd := GetPwd()
	messages := r.CreateConsumer()

	go func() {
		r.RabCtxDefine()
		for message := range messages {
			str := message.Body

			err := json.Unmarshal(str, &Input)
			if err != nil {
				log.Println(">>>>>>>>>> ", err)
				continue
			}

			err = m.DownloadVideo(Input.MinioBucket, Input.FolderName)
			if err != nil {
				log.Println(">>>>>>>>>> ", err)
				continue
			}
			RunFfmpeg()
			files, amountOfFiles := GetDirInfo()

			Output := Data{goalBucket, folder + strconv.Itoa(counter)}

			err = m.CreateNewBucket(Output.MinioBucket)
			if err != nil {
				log.Println(">>>>>>>>>> ", err)
				continue
			}
			ImageCompare(amountOfFiles, files, m.client, m.ctx, Output.MinioBucket, Output.FolderName)

			data, err := json.Marshal(Output)
			if err != nil {
				log.Println(">>>>>>>>>> ", err)
				continue
			}

			r.SendMessage(data)
			clear(pwd)
			counter++
		}

	}()
	<-forever
}
