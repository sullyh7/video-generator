package main

import (
	"fmt"
	"log"
	"sully/vid-gen-v2/service"
)

func main() {
	vs := service.NewVerseService("temp/audio.mp3", "temp/image.jpg", "output.mp4")
	verse, err := vs.GetVerseAndDownloadFiles()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(verse.Translation)
}
