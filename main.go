package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

const _videoCaptureTime = 5 * time.Second

func main() {
	messages := make(chan int)
	var recording bool

	go func() { runCamera(&recording, messages) }()
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
			log.Println("Waiting for all images to be processed..", len(messages), recording)
			for len(messages) > 0 || recording {
				log.Println("Waiting for all images to be processed..")
				time.Sleep(1 * time.Second)

				if counter > 10 {
					log.Println("Timeout reached, without all images being processed. Exiting...")
					break
				}
			}

			kill <- true
		}
	}()
	<-kill
}

func runCamera(recording *bool, messages chan int) {
	for {
		*recording = true
		log.Println("Video capture started...")

		_, err := exec.Command("libcamera-vid", "-t", "1000", "-o", "test.h264", "--width", "1920", "--height", "1080").Output()
		if err != nil {
			log.Fatalln(err)
		}

		messages <- 1
		*recording = false
		log.Println("Video captured successfully")
	}
}

func processVideo(messages chan int) {
	for {
		log.Println("Processing video...", len(messages))

		<-messages

		now := time.Now()
		// format now time to timestamp
		timestamp := now.Format("2006-01-02_15-04-05")

		_, err := exec.Command("ffmpeg", "-i", "test.h264", "-c:v", "copy", "-c:a", "copy", fmt.Sprint("html/", timestamp, ".mp4")).Output()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Processing video successfully completed")
	}
}
