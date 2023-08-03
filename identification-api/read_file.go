package main

import (
	"context"
	"encoding/base64"
	"io"
	"log"

	"cloud.google.com/go/storage"
)

func readFile() string {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	rc, err := client.Bucket("license-plates-go-parking").Object("testing-algo/frames11_output_0045.jpeg").NewReader(ctx)
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
