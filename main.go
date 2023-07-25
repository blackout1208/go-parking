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
	recording := make(chan int)
	process := make(chan int)

	go func() { runCamera(recording, process) }()
	go func() { processVideo(process) }()

	log.Println("Program started successfully..")

	c := make(chan os.Signal, 1)
	kill := make(chan bool)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Printf("captured %+v, exiting..", sig)
			// sig is a ^C, handle it

			counter := 0
			log.Println("Waiting for all images to be processed..", len(process))
			for len(process) > 0 {
				log.Println("Waiting for all images to be processed..")
				time.Sleep(_videoCaptureTime)

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

func runCamera(recording, process chan int) {
	for {
		log.Println("Video capture started...")
		recording <- 1
		_, err := exec.Command("libcamera-vid", "-t", fmt.Sprint(_videoCaptureTime), "-o", "test.h264", "--width", "1920", "--height", "1080").Output()
		if err != nil {
			log.Fatalln(err)
		}

		process <- 1
		<-recording
		log.Println("Video captured successfully")
	}
}

func processVideo(process chan int) {
	for {

		log.Println("Processing video...", len(process))
		<-process

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
