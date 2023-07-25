package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("captured %+v, exiting..", sig)
			// sig is a ^C, handle it
			touch()
		}
	}()

	messages := make(chan int)
	go func() { runCamera(messages) }()
	go func() { processVideo(messages) }()

	for {
		log.Println("Tracking started successfully")
	}
}

func runCamera(messages chan int) {
	for {
		_, err := exec.Command("libcamera-vid", "-t", "10000", "-o", "test.h264", "--width", "1920", "--height", "1080").Output()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Video captured successfully")

		messages <- 1
	}
}

func touch() {
	_, err := exec.Command("touch", "./flag.txt").Output()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("File touched successfully")
}

func processVideo(messages chan int) {
	for {
		<-messages

		now := time.Now()
		// format now time to timestamp
		timestamp := now.Format("2006-01-02_15-04-05")

		_, err := exec.Command("ffmpeg", "-i", "test.h264", "-c:v", "copy", "-c:a", "copy", fmt.Sprint("html/", timestamp, ".mp4")).Output()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Video processed successfully")
	}
}
