package gstorage

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

const scope = storage.DevstorageFullControlScope

func DownloadBucket(ctx context.Context, bucket, localDir string) error {
	client, err := google.DefaultClient(ctx, scope)
	if err != nil {
		return fmt.Errorf("failed to get client: %v", err)
	}
	service, err := storage.New(client)
	if err != nil {
		return fmt.Errorf("failed to get service: %v", err)
	}
	objs, err := service.Objects.List(bucket).Do()
	if err != nil {
		return fmt.Errorf("failed to list %v: %v", bucket, err)
	}
	for _, item := range objs.Items {
		fmt.Printf("%v\n", item.Name)
		downloadFile(service, bucket, item.Name, localDir)
	}
	return nil
}

func downloadFile(service *storage.Service, bucket, filename, dstDir string) error {
	resp, err := service.Objects.Get(bucket, filename).Download()
	if err != nil {
		return fmt.Errorf("failed to get file %q %q: %v", bucket, filename, err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath.Join(dstDir, filename))
	defer out.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error downloading %q: %v", filename, err)
	}

	log.Printf("Downloaded %v, %d bytes", filename, n)
	return nil
}
