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

type bqclient struct {
	service *bigquery.Service

	projectId string
	datasetId string

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

func NewBQClient() (*bqclient, error) {
	c := &bqclient{
		JobStatusPollingMaxTries: 40,
		JobStatusPollingInterval: 5 * time.Second,
	}
	service, err := getBigqueryService()
	if err != nil {
		return nil, err
	}
	c.service = service
	return c, nil
}

func (c *bqclient) JobQuery(query, dstTable string) *bigquery.Job {
	dstTableRef := &bigquery.TableReference{
		ProjectId: c.projectId,
		DatasetId: c.datasetId,
		TableId:   dstTable,
	}
	defaultDatasetRef := &bigquery.DatasetReference{
		ProjectId: c.projectId,
		DatasetId: c.datasetId,
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

func (c *bqclient) JobExtract(table, gspath string) *bigquery.Job {
	tableRef := &bigquery.TableReference{
		ProjectId: c.projectId,
		DatasetId: c.datasetId,
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
func (c *bqclient) RunSequentialJobs(jobs []*bigquery.Job) error {
	for _, job := range jobs {
		if err := c.RunJob(job); err != nil {
			return err
		}
	}
	return nil
}

// RunJob runs a bq Job synchronously.
func (c *bqclient) RunJob(job *bigquery.Job) error {
	job, err := c.service.Jobs.Insert(c.projectId, job).Do()
	if err != nil {
		return fmt.Errorf("failed to insert job: %v", err)
	}

	jobId := job.JobReference.JobId
	log.Printf("[Job %s] start polling ....", jobId)
	for i := 0; i < c.JobStatusPollingMaxTries; i++ {
		time.Sleep(c.JobStatusPollingInterval)
		j, err := c.service.Jobs.Get(c.projectId, jobId).Do()
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
