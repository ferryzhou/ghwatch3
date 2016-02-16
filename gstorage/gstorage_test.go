package gstorage

import (
	"testing"

	"golang.org/x/net/context"
)

func TestDownload(t *testing.T) {
	ctx := context.Background()
	if err := DownloadBucket(ctx, "bqclient_test", ""); err != nil {
		t.Errorf("DownloadBucket failed: %v, want nil", err)
	}
}
