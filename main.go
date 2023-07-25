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
	messages := make(chan int)
	recording := make(chan int)

	go func() { runCamera(messages, recording) }()
	go func() { processVideo(messages) }()

	log.Println("Program started successfully..")

	c := make(chan os.Signal, 1)
	kill := make(chan bool)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Printf("captured %+v, exiting..", sig)
			// sig is a ^C, handle it

			counter := 0
			for len(recording) > 0 {
				log.Println("Waiting for all images to be processed..")
				time.Sleep(1 * time.Second)

				if counter > 10 {
					log.Println("Timeout reached, without all images being processed. Exiting...")
					break
				}
			}

			log.Println("Killing libcamera-vid...")
			if _, err := exec.Command("pkill", "libcamera-vid").Output(); err != nil {
				log.Fatalln(err)
			}

			log.Println("Killing ffmpeg...")
			if _, err := exec.Command("pkill", "ffmpeg").Output(); err != nil {
				log.Fatalln(err)
			}

			kill <- true
		}
	}()
	<-kill
}

func runCamera(messages, recording chan int) {
	for {
		recording <- 1

		_, err := exec.Command("libcamera-vid", "-t", "10000", "-o", "test.h264", "--width", "1920", "--height", "1080").Output()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Video captured successfully")

		<-recording

		messages <- 1
	}
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
