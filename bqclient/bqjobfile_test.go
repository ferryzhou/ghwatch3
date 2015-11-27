package bqclient

import (
	"encoding/json"
	"testing"

	"google.golang.org/api/bigquery/v2"
)

func TestJobFile(t *testing.T) {
	c := newTestClient(t)
	testCases := []*struct {
		file    string
		wantJob *bigquery.Job
	}{
		{
			"testdata/01_query_test001.yml",
			c.JobQuery("SELECT * FROM [publicdata:samples.shakespeare] LIMIT 10", "test001"),
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
		jj, _ := json.Marshal(job)
		wj, _ := json.Marshal(tc.wantJob)
		if string(jj) != string(wj) {
			t.Errorf("got job \n%v, want \n%v", string(jj), string(wj))
		}
	}
}
