// Program pipeline run bq queries and download data to local.
// BQ jobs are configured in jobs/ folder.
// Jobs are executed sequentially.
// All files in cloud storage in specified bucket are downloaded to local after all jobs finished.
// Example:
//   go run pipeline.go --jobs_path=jobs/01_*.yml  --project=950350008903 --dst_dir=results
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	bq "github.com/ferryzhou/ghwatch3/bqclient"
	gs "github.com/ferryzhou/ghwatch3/gstorage"
)

// This name is used in sql specified in jobs/ folder.
const bucket = "ghwatch3"
const dataset = "ghwatch3"

var (
	project  = flag.String("project", "", "bigquery project id")
	jobsPath = flag.String("jobs_path", "", "job files pattern, e.g. jobs/*.yml")
	dstDir   = flag.String("dst_dir", ".", "destination folder storing the result files")
)

func newBQClient() *bq.BQClient {
	c, err := bq.NewBQClient()
	if err != nil {
		log.Panicf("failed to start client: %v", err)
	}
	c.ProjectId = *project
	c.DatasetId = dataset
	return c
}

func main() {
	flag.Parse()
	fmt.Printf("job dir: %v\n", *jobsPath)
	c := newBQClient()
	if err := c.RunJobsInFolder(*jobsPath); err != nil {
		log.Panicf("Failed to run jobs: %v", err)
	}
	ctx := context.Background()
	if err := os.MkdirAll(*dstDir, 0777); err != nil {
		log.Panicf("Failed to create dst dir %v: %v", *dstDir, err)
	}
	if err := gs.DownloadBucket(ctx, bucket, *dstDir); err != nil {
		log.Panicf("Failed to download bucket %v: %v", bucket, err)
	}
	fmt.Printf("Success")
}
