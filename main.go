package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"time"
)

const _videoCaptureTime = 5 * time.Second

var wg sync.WaitGroup

func main() {

	go func() { runCamera() }()

	log.Println("Program started successfully..")

	c := make(chan os.Signal, 1)
	kill := make(chan bool)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Println("captured ^C, exiting..", sig)

			log.Println("Waiting for all images to be processed..")
			wg.Wait()

			kill <- true
		}
	}()
	<-kill
}

func runCamera() {
	for {
		log.Println("Video capture started...")

		_, err := exec.Command("libcamera-vid", "-t", fmt.Sprint(_videoCaptureTime), "-o", "test.h264", "--width", "1920", "--height", "1080").Output()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Video captured successfully")

		wg.Add(1)
		go func() {
			defer wg.Done()
			processVideo()
		}()
	}
}

func processVideo() {
	log.Println("Processing video...")

	now := time.Now()
	// format now time to timestamp
	timestamp := now.Format("2006-01-02_15-04-05")

	_, err := exec.Command("ffmpeg", "-i", "test.h264", "-c:v", "copy", "-c:a", "copy", fmt.Sprint("html/", timestamp, ".mp4")).Output()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Processing video successfully completed")
}
