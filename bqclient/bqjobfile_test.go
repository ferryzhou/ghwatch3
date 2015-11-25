package bqclient

import (
	"reflect"
	"testing"

	"google.golang.org/cloud/bigquery"
)

func TestJobFile(t *testing.T) {
	pemPath := "g.pem"
	c, err := NewBQClient(pemPath)
	if err != nil {
		t.Errorf("failed to start client: %v", err)
	}

	testCases := []*struct {
		file    string
		wantJob *bigquery.Job
	}{
		{
			"testdata/01_query_test001.yml",
			c.JobQuery("test001", "SELECT * FROM [publicdata:samples.shakespeare] LIMIT 10"),
		},
		{
			"testdata/02_extract_test001.yml",
			c.JobExtract("test001", "gs://bqclient_test/bqclient_test.gzip"),
		},
	}

	for _, tc := range testCases {
		job, err := c.JobFromFile(tc.file)
		if err != nil {
			t.Errorf("failed to get job from file: %v", err)
			continue
		}
		if !reflect.DeepEqual(job, tc.wantJob) {
			t.Errorf("got job %v, want %v", job, tc.wantJob)
		}
	}
}
