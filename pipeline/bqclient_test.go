package pipeline

import "testing"

func TestBqJobs(t *testing.T) {
	c, err := NewBQClient(pemPath)
	if err != nil {
		t.Errorf("failed to start client: %v", err)
	}
	// Configure a serial job pipeline.
	// first query to table, second query to table, and then an extract.

}
