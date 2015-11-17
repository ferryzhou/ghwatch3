package pipeline

import (
	"testing"

	"google.golang.org/api/bigquery/v2"
)

const (
	test_project = "950350008903"
	test_dataset = "bqclient_test"
	test_bucket  = "bqclient_test"
)

// Before testing, make sure
//   1) pem file is copied to g.pem
//   2) a dataset called bqclient_test is created.
//   3) a data storage bucket bqclient_test is created.
func TestBqJobs(t *testing.T) {
	pemPath := "g.pem"
	c, err := NewBQClient(pemPath)
	if err != nil {
		t.Errorf("failed to start client: %v", err)
	}
	c.projectId = test_project
	c.datasetId = test_dataset
	// Configure a serial job pipeline.
	// first query to table, second query to table, and then an extract.
	jobs := []*bigquery.Job{
		c.JobQuery("SELECT * FROM [publicdata:samples.shakespeare] LIMIT 10", "bqclient_test_1"),
		c.JobQuery("SELECT word, word_count, FROM bqclient_test_1", "bqclient_test_2"),
		c.JobExtract("bqclient_test_2", "gs://"+test_bucket+"/bqclient_test.gzip"),
	}
	if err := c.RunSequentialJobs(jobs); err != nil {
		t.Errorf("faield to run jobs: %v", err)
	}
}
