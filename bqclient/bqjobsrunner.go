package bqclient

import (
	"fmt"
	"io/ioutil"
	"path"

	"google.golang.org/api/bigquery/v2"
)

// RunJobsInFolder scan jobs file in a folder and run them sequentially.
// Files are sorted by filename.
// Subdirectories are not scanned.
func (c *bqclient) RunJobsInFolder(dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return fmt.Errorf("failed to read dir %v: %v", dirname, err)
	}
	var jobs []*bigquery.Job
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		job, err := c.JobFromFile(path.Join(dirname, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read file: %v", err)
		}
		jobs = append(jobs, job)
	}
	if err := c.RunSequentialJobs(jobs); err != nil {
		return fmt.Errorf("failed to run jobs: %v", err)
	}
	return nil
}
