package main

import (
	"context"
	"encoding/base64"
	"io"
	"log"

	"cloud.google.com/go/storage"
)

func GetFrame(bucketName, frameURI string) string {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	rc, err := client.Bucket(bucketName).Object(frameURI).NewReader(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	body, err := io.ReadAll(rc)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(body)
}
