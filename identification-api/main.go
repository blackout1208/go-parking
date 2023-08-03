package main

import (
	"log"

	"github.com/joho/godotenv"
)

type Test struct {
	MimeType string `protobuf:"mimeType"`
	Content  string `protobuf:"content"`
}

func init() {
	loadSecrets()
}

func main() {
	frameBase64 := readFile()

	response := predictLicensePlate(frameBase64)
	prediction := extractPrediction(response)

	var indexes []int
	for i, confidence := range prediction.Confidences {
		if confidence < 0.2 || prediction.DisplayNames[i] != _licensePlateLabel {
			continue
		}

		indexes = append(indexes, i)
	}
	// "bboxes": [ [xMin, xMax, yMin, yMax], ...]
	var images []bbox
	for _, index := range indexes {
		images = append(images, bbox{
			xmin: prediction.Bboxes[index][0],
			xmax: prediction.Bboxes[index][1],
			ymin: prediction.Bboxes[index][2],
			ymax: prediction.Bboxes[index][3],
		})
	}

	for i, image := range images {
		// extract license plate from video frame
		framePlate := extractLicensePlateIMG(i, frameBase64, bbox{
			xmin: image.xmin,
			xmax: image.xmax,
			ymin: image.ymin,
			ymax: image.ymax,
		})

		// convert image with license to text
		licenses := findTextInFrame(framePlate)

		for _, license := range licenses {
			// send license to database
			log.Printf("License plate: %+v\n", license)
		}
	}
}

func loadSecrets() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Error loading .env file, continuing...", err)
	}
}
