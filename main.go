package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
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
		inputChannel <- input
		if input == '0' {
			fmt.Println("Starting..")
			go runJPEG(inputChannel)
		} else if input == '1' {
			fmt.Println("Starting..")
			go runCamera(inputChannel)
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

func runJPEG(inputChannel chan rune) {
	for {
		select {
		case input, ok := <-inputChannel:
			if ok && input != '0' {
				fmt.Println("Stopping..")
				return
			} else {
				fmt.Println("Channel closed!", input, ok)
				return
			}
		default:
		}

		now := time.Now()
		// format now time to timestamp
		timestamp := now.Format("2006-01-02_15-04-05")

		_, err := exec.Command("libcamera-jpeg", "-o", fmt.Sprint("./", "tmp/", timestamp, ".jpg"), "--width", "1920", "--height", "1080").Output()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Picture captured successfully")

	}
}

func runCamera(inputChannel chan rune) {
	for {
		select {
		case input, ok := <-inputChannel:
			if ok && input != '1' {
				fmt.Println("Stopping..")
				return
			} else {
				fmt.Println("Channel closed!")
				return
			}
		default:
		}

		fmt.Println("Video capture started...")

		now := time.Now()
		// format now time to timestamp
		timestamp := now.Format("2006-01-02_15-04-05")

		_, err := exec.Command("libcamera-vid", "-t", "1000", "-o", fmt.Sprint("./", "tmp/", timestamp, ".h264"), "--width", "1920", "--height", "1080").Output()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Video captured successfully")
	}
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

		_, err := exec.Command("ffmpeg", "-i", fmt.Sprint("./", "tmp/", file.Name()), "-c:v", "copy", "-c:a", "copy", fmt.Sprint("./", "html/", file.Name(), ".mp4")).Output()
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Println("Processing video successfully completed")
}
