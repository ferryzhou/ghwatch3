package bqclient

import (
	"fmt"
	"io/ioutil"

	"google.golang.org/api/bigquery/v2"
	yaml "gopkg.in/yaml.v2"
)

type jobConfig struct {
	Command string
	Table   string
	Gspath  string
	Query   string
}

// parseJobFile generate a bigQuery job by parsing a yaml file.
func (c *bqclient) JobFromFile(file string) (*bigquery.Job, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %v: %v", file, err)
	}
	var jc jobConfig
	if err := yaml.Unmarshal(data, &jc); err != nil {
		return nil, fmt.Errorf("failed to parse file %v: %v", file, err)
	}
	fmt.Printf("%v", string(data))
	fmt.Printf("%v", jc)
	switch jc.Command {
	case "query":
		return c.JobQuery(jc.Query, jc.Table), nil
	case "extract":
		return c.JobExtract(jc.Table, jc.Gspath), nil
	default:
		return nil, fmt.Errorf("no such Command %v", jc.Command)
	}
}
