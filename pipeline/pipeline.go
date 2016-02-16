// Program pipeline run bq queries and download data to local.
// BQ jobs are configured in jobs/ folder.
// Jobs are executed sequentially.
// All files in cloud storage in specified bucket are downloaded to local after all jobs finished.
package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/net/context"

	bq "github.com/ferryzhou/ghwatch3/bqclient"
	gs "github.com/ferryzhou/ghwatch3/gstorage"
)

// This name is used in sql specified in jobs/ folder.
const bucket = "ghwatch3"

var (
	project = flag.String("project", "", "bigquery project id")
	dataset = flag.String("dataset", "", "dest bigquery dataset name")
	jobsDir = flag.String("jobs_dir", "", "a folder containing job files")
)

func newBQClient() *bq.BQClient {
	c, err := bq.NewBQClient()
	if err != nil {
		log.Panicf("failed to start client: %v", err)
	}
	c.ProjectId = *project
	c.DatasetId = *dataset
	return c
}

func main() {
	flag.Parse()
	fmt.Printf("job dir: %v\n", *jobsDir)
	c := newBQClient()
	if err := c.RunJobsInFolder(*jobsDir); err != nil {
		log.Panicf("failed to run jobs: %v", err)
	}
	ctx := context.Background()
	if err := gs.DownloadBucket(ctx, bucket, "."); err != nil {
		log.Panicf("failed to download bucket %v: %v", bucket, err)
	}
	fmt.Printf("Success")
}
