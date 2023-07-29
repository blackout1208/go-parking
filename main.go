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

var wg sync.WaitGroup

func main() {
	kill := make(chan bool)

	go func() { runCamera(kill) }()

	log.Println("Program started successfully..")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Println("captured ^C, exiting..", sig)

			log.Println("Waiting for all images to be processed..")
			wg.Wait()
			log.Println("All images processed..")

			kill <- true
		}
	}()
	<-kill
}

func runCamera(kill chan bool) {
	for {
		select {
		case isToStop := <-kill:
			kill <- isToStop
			return
		default:
		}

		log.Println("Video capture started...")
		wg.Add(1)

		_, err := exec.Command("libcamera-vid", "-t", "1000", "-o", "test.h264", "--width", "1920", "--height", "1080").Output()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Video captured successfully")

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
