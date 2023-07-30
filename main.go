package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

var wg sync.WaitGroup

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
	err := keyboard.Open()
	if err != nil {
		log.Fatalln("Error opening keyboard:", err)
		return
	}
	defer keyboard.Close()

	inputChannel := make(chan rune, 1)

	for {
		readInput(&inputChannel)
		input := <-inputChannel
		if input == '1' {
			fmt.Println("Starting..")
			go runCamera()
		} else if input == '2' {
			fmt.Println("Exiting..")
			wg.Wait()
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
		wg.Add(1)

		_, err := exec.Command("libcamera-vid", "-t", "1000", "-o", "test.h264", "--width", "1920", "--height", "1080").Output()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Video captured successfully")

		go func() {
			defer wg.Done()
			processVideo()
		}()
	}
}

func processVideo() {
	fmt.Println("Processing video...")

	now := time.Now()
	// format now time to timestamp
	timestamp := now.Format("2006-01-02_15-04-05")

	_, err := exec.Command("ffmpeg", "-i", "test.h264", "-c:v", "copy", "-c:a", "copy", fmt.Sprint("html/", timestamp, ".mp4")).Output()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Processing video successfully completed")
}
