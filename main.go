package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func main() {
	inputChannel := make(chan rune, 1)
	runtime := time.Now().Format("2006-01-02_15-04-05")

	for {
		readInput(&inputChannel)
		input := <-inputChannel
		inputChannel <- input
		if strings.ContainsRune("1", input) {
			fmt.Println("Starting..")

			runtime = time.Now().Format("2006-01-02_15-04-05")
			rootFolder := fmt.Sprint("./", runtime, "/")

			go runCamera(rootFolder, inputChannel)
		} else if strings.ContainsRune("2", input) {
			fmt.Println("Exiting..")

			rootFolder := fmt.Sprint("./", runtime, "/")

			processVideo(rootFolder)
			os.Exit(0)
			break
		} else {
			fmt.Println("Invalid option selected")
		}
	}
}

func readInput(inputChannel *chan rune) {
	fmt.Println(`Choose an option: 1. Start; 2. Exit`)

	reader := bufio.NewReader(os.Stdin)
	input, _, err := reader.ReadRune()
	if err != nil {
		log.Fatalln(err)
	}

	*inputChannel <- input
}

func runCamera(rootFolder string, inputChannel chan rune) {
	select {
	case input, ok := <-inputChannel:
		if ok && !strings.ContainsRune("1", input) {
			fmt.Println("Stopping..")

			output, err := exec.Command("pkill", "libcamera-vid").Output()
			if err != nil {
				log.Fatalln(output, err)
			}

			return
		} else if strings.ContainsRune("2", input) {
			fmt.Println("Channel closed!")
			return
		}
	default:
	}

	if err := createDirectories(fmt.Sprint(rootFolder)); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Video capture started...")

	fileName := time.Now().Format("2006-01-02_15-04-05")

	_, err := exec.Command("libcamera-vid", "-t", "0", "-o", fmt.Sprint(rootFolder, fileName, ".h264"), "--width", "1920", "--height", "1080").Output()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Video captured successfully")
}

func processVideo(rootFolder string) {
	fmt.Println("Processing videos...")

	files, err := os.ReadDir(rootFolder)
	if err != nil {
		log.Fatal(err)
	}

	mp4Folder := fmt.Sprint(rootFolder, "mp4/")
	framesFolder := fmt.Sprint(rootFolder, "frames/")

	createDirectories(mp4Folder)
	createDirectories(framesFolder)

	for _, file := range files {
		if path.Ext(file.Name()) != ".h264" {
			continue
		}

		_, err := exec.Command("ffmpeg", "-i", fmt.Sprint(rootFolder, file.Name()), "-c:v", "copy", "-c:a", "copy", fmt.Sprint(mp4Folder, file.Name(), ".mp4")).Output()
		if err != nil {
			log.Fatalln(err)
		}

		_, err = exec.Command("ffmpeg", "-i", fmt.Sprint(mp4Folder, file.Name(), ".mp4"), "-r", "1", fmt.Sprint(framesFolder, file.Name(), "_%04d.png")).Output()
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Println("Processing video successfully completed")
}

func createDirectories(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
