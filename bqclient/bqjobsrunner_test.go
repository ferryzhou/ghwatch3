package bqclient

import (
	"testing"
)

func TestJobsRunner(t *testing.T) {
	c := newTestClient(t)
	testCases := []*struct {
		dirname string
		wantErr bool
	}{
		{
			"testdata/goodjobs01/*.yml",
			false,
		},
		{
			"testdata/badjobs01/*.yml",
			true,
		},
	}

	for _, tc := range testCases {
		if err := c.RunJobsInFolder(tc.dirname); tc.wantErr != (err != nil) {
			t.Errorf("got err: %v, wantErr: %v", err, tc.wantErr)
		}
	}
}
