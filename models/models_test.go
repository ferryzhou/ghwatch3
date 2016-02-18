package models

import (
	"reflect"
	"testing"
)

func TestInMemRepoStore(t *testing.T) {
	repoCSV := "testdata/repos.csv"
	recCSV := ""
	s, err := InMemRepoStoreFromCSV(repoCSV, recCSV)
	if err != nil {
		t.Fatalf("InMemRepoStoreFromCSV(%q, %q) failed: %v", repoCSV, recCSV, err)
	}
	want := &Repo{
		ShortPath:    "joyent/node",
		URL:          "https://github.com/joyent/node",
		Stars:        33984,
		Language:     "JavaScript",
		Description:  "evented I/O for v8 javascript",
		name:         "node",
		owner:        "joyent",
		organization: "joyent",
		createdAt:    "2009-05-27 16:29:46",
		pushedAt:     "2014-12-31 21:16:33",
	}
	if got := s.RepoByShortPath(want.ShortPath); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
