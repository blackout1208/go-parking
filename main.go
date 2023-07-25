package main

import (
	"log"
	"os/exec"
)

func main() {
	runCamera()
}

func runCamera() {
	_, err := exec.Command("libcamera-vid", "-t", "10000", "-o", "test.h264", "--width", "1920", "--height", "1080").Output()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Video captured successfully")

	_, err = exec.Command("touch", "./flag.txt").Output()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("File touched successfully")
}
