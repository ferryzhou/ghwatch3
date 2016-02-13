package main

import (
	"flag"
	"fmt"
	"log"

	bq "github.com/ferryzhou/ghwatch3/bqclient"
)

var (
	jobsDir = flag.String("jobs_dir", "", "a folder containing job files")
	project = flag.String("project", "", "bigquery project id")
	dataset = flag.String("dataset", "", "bigquery dataset name")
	bucket  = flag.String("bucket", "", "google cloud storage bucket name")
)

func newClient() *bq.BQClient {
	c, err := bq.NewBQClient()
	if err != nil {
		log.Panicf("failed to start client: %v", err)
	}
	c.ProjectId = *project
	c.DatasetId = *dataset
	return c
}

func main() {
	fmt.Printf("job dir: %v", *jobsDir)
	c := newClient()
	if err := c.RunJobsInFolder(*jobsDir); err != nil {
		log.Panicf("failed to run jobs: %v", err)
	}
	fmt.Printf("Success")
}
