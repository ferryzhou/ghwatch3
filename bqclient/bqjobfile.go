package bqclient

import (
	"google.golang.org/api/bigquery/v2"
)

// parseJobFile generate a bigquery job by parsing a yaml file.
func (c *bqclient) JobFromFile(file string) (*bigquery.Job, error) {
	return nil, nil
}
