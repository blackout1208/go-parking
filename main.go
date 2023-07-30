package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
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
	inputChannel := make(chan rune, 1)

	for {
		readInput(&inputChannel)
		input := <-inputChannel
		inputChannel <- input
		if strings.ContainsRune("1", input) {
			fmt.Println("Starting..")
			go runCamera(inputChannel)
		} else if strings.ContainsRune("2", input) {
			fmt.Println("Exiting..")
			processVideo()
			os.Exit(0)
			break
		} else {
			fmt.Println("Invalid option selected")
		}
	}
}

func runCamera(inputChannel chan rune) {
	select {
	case input, ok := <-inputChannel:
		if ok && !strings.ContainsRune("1", input) {
			fmt.Println("Stopping..")

			_, err := exec.Command("pkill", "libcamera-vid").Output()
			if err != nil {
				log.Fatalln(err)
			}

			return
		} else if strings.ContainsRune("2", input) {
			fmt.Println("Channel closed!")
			return
		}
	default:
	}

	fmt.Println("Video capture started...")

	now := time.Now()
	timestamp := now.Format("2006-01-02_15-04-05")

	_, err := exec.Command("libcamera-vid", "-t", "0", "-o", fmt.Sprint("./", "tmp/", timestamp, ".h264"), "--width", "1920", "--height", "1080").Output()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Video captured successfully")
}

func processVideo() {
	fmt.Println("Processing videos...")

	files, err := os.ReadDir(fmt.Sprint("./", "tmp/"))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if path.Ext(file.Name()) != ".h264" {
			continue
		}

		_, err := exec.Command("ffmpeg", "-i", fmt.Sprint("./", "tmp/", file.Name()), "-c:v", "copy", "-c:a", "copy", fmt.Sprint("./", "mp4/", file.Name(), ".mp4")).Output()
		if err != nil {
			log.Fatalln(err)
		}

		_, err = exec.Command("ffmpeg", "-i", fmt.Sprint("./", "mp4/", file.Name(), ".mp4"), "-r", "1", fmt.Sprint("./", "frames/", file.Name(), "_%04d.png")).Output()
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Println("Processing video successfully completed")
}
