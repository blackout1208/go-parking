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

	// Get frame from cloud storage
	bucketName := "license-plates-go-parking"
	fileURI := "testing-algo/frames9_output_0015.jpeg"
	frameBase64 := GetFrame(bucketName, fileURI)

	response := predictLicensePlate(frameBase64)
	prediction := extractPrediction(response)

	platesIMG := prediction.getPlatesIMG()

	for i, image := range platesIMG {
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
