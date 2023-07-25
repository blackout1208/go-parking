package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
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

	runCamera()
}

func runCamera() {
	_, err := exec.Command("libcamera-vid", "-t", "10000", "-o", "test.h264", "--width", "1920", "--height", "1080").Output()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Video captured successfully")

	touch()
}

func touch() {
	_, err := exec.Command("touch", "./flag.txt").Output()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("File touched successfully")
}
