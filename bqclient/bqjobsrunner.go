package bqclient

import (
	"fmt"
	"path/filepath"

	"google.golang.org/api/bigquery/v2"
)

// RunJobsInFolder scan jobs file in a folder and run them sequentially.
// Files are sorted by filename.
// Subdirectories are not scanned.
func (c *BQClient) RunJobsInFolder(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to read dir %v: %v", pattern, err)
	}
	var jobs []*bigquery.Job
	for _, file := range files {
		job, err := c.JobFromFile(file)
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
