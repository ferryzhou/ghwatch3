// package pipeline contains tools to
//  1. move data from bq to local.
//  2. run queries on bq to extract data
//  3. store dumped data to database
package pipeline

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/oauth2"
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
}

func getBigqueryService(pemPath string) (*bigquery.Service, error) {
	// generate auth token and create service object
	pemKeyBytes, err := ioutil.ReadFile(pemPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pem file: %s, error: %v", pemPath, err)
	}

	conf, err := google.JWTConfigFromJSON(pemKeyBytes, bqEndpoint)

	client := conf.Client(oauth2.NoContext)

	return bigquery.New(client)
}

func NewBQClient(pemPath string) (*bqclient, error) {
	c := &bqclient{}
	service, err := getBigqueryService(pemPath)
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
		FlattenResults:    flatten,
	}
	conf := &bigquery.JobConfiguration{
		Query: qConf,
	}

	return &bigquery.Job{
		Configuration: conf,
	}
}

func (c *bqclient) JobExtrac(table, gspath string) *bigquery.Job {
	tableRef := &bigquery.TableReference{
		ProjectId: c.projectId,
		DatasetId: c.datasetId,
		TableId:   table,
	}
	extract := &bigquery.JobConfigurationExtract{
		SourceTable:       tableRef,
		DestinationUris:   []string{*gspath},
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

func (c *bqclient) startJob(j *bigquery.Job) (*bigquery.Job, error) {
	return service.Jobs.Insert(c.projectId, job).Do()
}
