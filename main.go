package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func readInput(inputChannel *chan rune) {
	fmt.Println(`Choose an option: 1. Start; 2. Exit`)

	reader := bufio.NewReader(os.Stdin)
	input, _, err := reader.ReadRune()
	if err != nil {
		log.Fatalln(err)
	}

	*inputChannel <- input
}

func main() {
	// err := keyboard.Open()
	// if err != nil {
	// 	log.Fatalln("Error opening keyboard:", err)
	// 	return
	// }
	// defer keyboard.Close()

	inputChannel := make(chan rune, 1)

	for {
		readInput(&inputChannel)
		input := <-inputChannel
		if input == '1' {
			fmt.Println("Starting..")
			runCamera()
		} else if input == '2' {
			fmt.Println("Exiting..")
			processVideo()
			os.Exit(0)
			break
		} else {
			fmt.Println("Invalid option selected")
		}
	}
}

func runCamera() {
	for {
		fmt.Println("Video capture started...")

		now := time.Now()
		// format now time to timestamp
		timestamp := now.Format("2006-01-02_15-04-05")

		_, err := exec.Command("libcamera-vid", "-t", "1000", "-o", fmt.Sprint("/tmp/", timestamp, ".h264"), "--width", "1920", "--height", "1080").Output()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Video captured successfully")

		// go processVideo()
	}
}

func processVideo() {
	fmt.Println("Processing videos...")

	files, err := os.ReadDir("/tmp/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fmt.Println(file.Name(), file.IsDir())

		now := time.Now()
		// format now time to timestamp
		timestamp := now.Format("2006-01-02_15-04-05")

		_, err := exec.Command("ffmpeg", "-i", fmt.Sprint("/tmp/", file, ".h264"), "-c:v", "copy", "-c:a", "copy", fmt.Sprint("html/", timestamp, ".mp4")).Output()
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Println("Processing video successfully completed")
}
