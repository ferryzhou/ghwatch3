package main

import (
	"flag"
	"io"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

const scope = storage.DevstorageFullControlScope

var (
	bucket   = flag.String("bucket", "", "gcloud storage bucket name")
	filename = flag.String("filename", "", "filename in gcloud storage bucket")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, scope)
	if err != nil {
		log.Fatalf("failed to get client: %v", err)
	}
	service, err := storage.New(client)

	resp, err := service.Objects.Get(*bucket, *filename).Download()
	if err != nil {
		log.Fatalf("failed to get file %q %q: %v", *bucket, *filename, err)
	}
	defer resp.Body.Close()

	out, err := os.Create(*filename)
	defer out.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("error downloading %q: %v", *filename, err)
	}

	log.Printf("Downloaded %d bytes", n)
}
