// package bigquery contains functions to
//  1. run bigquery jobs such as query to table and export table to google storage.
package bqclient

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/bigquery/v2"
)

const (
	bqEndpoint = "https://www.googleapis.com/auth/bigquery"
)

type BQClient struct {
	service *bigquery.Service

	ProjectId string
	DatasetId string

	JobStatusPollingInterval time.Duration
	JobStatusPollingMaxTries int
}

func getBigqueryService() (*bigquery.Service, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, bqEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}

	return bigquery.New(client)
}

func NewBQClient() (*BQClient, error) {
	c := &BQClient{
		JobStatusPollingMaxTries: 400,
		JobStatusPollingInterval: 5 * time.Second,
	}
	service, err := getBigqueryService()
	if err != nil {
		return nil, err
	}
	c.service = service
	return c, nil
}

func (c *BQClient) JobQuery(query, dstTable string) *bigquery.Job {
	dstTableRef := &bigquery.TableReference{
		ProjectId: c.ProjectId,
		DatasetId: c.DatasetId,
		TableId:   dstTable,
	}
	defaultDatasetRef := &bigquery.DatasetReference{
		ProjectId: c.ProjectId,
		DatasetId: c.DatasetId,
	}
	qConf := &bigquery.JobConfigurationQuery{
		Query:             query,
		DestinationTable:  dstTableRef,
		DefaultDataset:    defaultDatasetRef,
		AllowLargeResults: true,
		WriteDisposition:  "WRITE_TRUNCATE",
		CreateDisposition: "CREATE_IF_NEEDED",
	}
	conf := &bigquery.JobConfiguration{
		Query: qConf,
	}

	return &bigquery.Job{
		Configuration: conf,
	}
}

func (c *BQClient) JobExtract(table, gspath string) *bigquery.Job {
	tableRef := &bigquery.TableReference{
		ProjectId: c.ProjectId,
		DatasetId: c.DatasetId,
		TableId:   table,
	}
	extract := &bigquery.JobConfigurationExtract{
		SourceTable:       tableRef,
		DestinationUris:   []string{gspath},
		DestinationFormat: "CSV",
		Compression:       "GZIP",
	}
	conf := &bigquery.JobConfiguration{
		Extract: extract,
	}

	return &bigquery.Job{
		Configuration: conf,
	}
}

// run a series of jobs sequentially and synchronously
func (c *BQClient) RunSequentialJobs(jobs []*bigquery.Job) error {
	for _, job := range jobs {
		if err := c.RunJob(job); err != nil {
			return err
		}
	}
	return nil
}

// RunJob runs a bq Job synchronously.
func (c *BQClient) RunJob(job *bigquery.Job) error {
	log.Printf("Run job %v", job)
	job, err := c.service.Jobs.Insert(c.ProjectId, job).Do()
	if err != nil {
		return fmt.Errorf("failed to insert job: %v", err)
	}

	jobId := job.JobReference.JobId
	log.Printf("[Job %s] start polling ....", jobId)
	for i := 0; i < c.JobStatusPollingMaxTries; i++ {
		time.Sleep(c.JobStatusPollingInterval)
		j, err := c.service.Jobs.Get(c.ProjectId, jobId).Do()
		if err != nil {
			log.Printf("[Job %s] failed to get job status: %v\n", jobId, err)
			continue
		}
		log.Printf("[Job %s] status: %s\n", jobId, j.Status.State)
		if j.Status.State != "DONE" {
			continue
		}
		if err := j.Status.ErrorResult; err != nil {
			return fmt.Errorf("[Job %s] job failed: %v", err)
		}

		return nil
	}
	return fmt.Errorf("Timeout")
}
