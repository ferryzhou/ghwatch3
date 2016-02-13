package bqclient

import (
	"testing"

	"google.golang.org/api/bigquery/v2"
)

const (
	test_project = "950350008903"
	test_dataset = "bqclient_test"
	test_bucket  = "bqclient_test"
)

func newTestClient(t *testing.T) *BQClient {
	c, err := NewBQClient()
	if err != nil {
		t.Fatalf("failed to start client: %v", err)
	}
	c.ProjectId = test_project
	c.DatasetId = test_dataset
	return c
}

// Before testing, make sure
//   1) pem file is copied to g.pem
//   2) a dataset called BQClient_test is created.
//   3) a data storage bucket BQClient_test is created.
func TestBqJobs(t *testing.T) {
	c := newTestClient(t)
	// Configure a serial job pipeline.
	// first query to table, second query to table, and then an extract.
	jobs := []*bigquery.Job{
	//c.JobQuery("SELECT * FROM [publicdata:samples.shakespeare] LIMIT 10", "BQClient_test_1"),
	//c.JobQuery("SELECT word, word_count, FROM BQClient_test_1", "BQClient_test_2"),
	//		c.JobExtract("BQClient_test_2", "gs://"+test_bucket+"/BQClient_test.gzip"),
	}
	if err := c.RunSequentialJobs(jobs); err != nil {
		t.Errorf("faield to run jobs: %v", err)
	}
}
