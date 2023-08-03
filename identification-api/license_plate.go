package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	vision "cloud.google.com/go/vision/apiv1"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	_licensePlateConfidenceThreshold = 0.2
	_licensePlateMaxPredictions      = 5
)

// https://cloud.google.com/vertex-ai/docs/image-data/object-detection/interpret-results
// "bboxes": [ [xMin, xMax, yMin, yMax], ...]
type bbox struct {
	xmin float64
	ymin float64
	xmax float64
	ymax float64
}

func predictLicensePlate(frameBase64 string) map[string]interface{} {
	ctx := context.Background()

	c, err := aiplatform.NewPredictionClient(
		ctx,
		option.WithEndpoint("us-central1-aiplatform.googleapis.com:443"),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	parameters, err := structpb.NewValue(map[string]interface{}{
		"confidenceThreshold": _licensePlateConfidenceThreshold,
		"maxPredictions":      _licensePlateMaxPredictions,
	})
	if err != nil {
		log.Printf("QueryVertex structpb.NewValue parameters - Err:%s", err)
	}

	instance, err := structpb.NewValue(map[string]interface{}{
		"content": frameBase64,
	})
	if err != nil {
		log.Printf("QueryVertex structpb.NewValue instance - Err:%s", err)
	}

	// create aiplatformpb.predictRequest object with the necessary instances to send an image/jpeg to the endpoint
	req := &aiplatformpb.PredictRequest{
		Endpoint:   "projects/1086189634506/locations/us-central1/endpoints/1195419828043644928",
		Instances:  []*structpb.Value{instance},
		Parameters: parameters,
	}

	resp, err := c.Predict(ctx, req)
	if err != nil {
		log.Fatalln(err)
	}

	respMap := resp.Predictions[0].GetStructValue().AsMap()
	return respMap
}

// extractLicensePlate extracts the license plate from the image.
func extractLicensePlateIMG(i int, frame string, area bbox) string {
	fsrc := base64.NewDecoder(base64.StdEncoding, strings.NewReader(frame))

	src, err := jpeg.Decode(fsrc)
	if err != nil {
		log.Fatal(err)
	}

	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y

	minX := int(area.xmin * float64(width))
	minY := int(area.ymin * float64(height))
	maxX := int(area.xmax * float64(width))
	maxY := int(area.ymax * float64(height))

	dst := image.NewRGBA(image.Rect(minX, minY, maxX, maxY))

	draw.Draw(dst, image.Rect(minX, minY, maxX, maxY), src, image.Point{minX, minY}, draw.Src)

	file, err := os.Create(fmt.Sprint("draw", i, ".png"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, dst); err != nil {
		log.Fatal(err)
	}

	// create a base64 encoding of the cropped image
	var buf strings.Builder
	if err := png.Encode(&buf, dst); err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString([]byte(buf.String()))
}

func findTextInFrame(frame string) []string {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	fsrc := base64.NewDecoder(base64.StdEncoding, strings.NewReader(frame))

	image, err := vision.NewImageFromReader(fsrc)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	labels, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	var licenses []string
	for _, label := range labels {
		if len(label.Description) <= 2 {
			continue
		}

		licenses = append(licenses, label.Description)
	}

	return licenses
}
