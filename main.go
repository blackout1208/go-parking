package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"golang.org/x/term"
)

var wg sync.WaitGroup

func readInput(inputChannel *chan rune) {
	fmt.Print(`Choose an option: 1. Start; 2. Exit`)

	reader := bufio.NewReader(os.Stdin)
	input, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}

	*inputChannel <- input
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

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
			log.Fatal(err)
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
		log.Fatal(err)
	}

	fmt.Println("Processing video successfully completed")
}
