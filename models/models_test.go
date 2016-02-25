package models

import (
	"reflect"
	"testing"
)

const (
	reposCSV     = "testdata/repos.csv"
	reposCSVGZip = "testdata/repos.csv.gz"
	recsCSVGZip  = "testdata/recs.csv.gz"
)

func TestInMemRepoStore(t *testing.T) {
	s := NewInMemRepoStore()
	if err := s.LoadReposFromCSVGZip(reposCSVGZip); err != nil {
		t.Fatalf("LoadReposFromCSVGZip(%q) failed: %v", reposCSVGZip, err)
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
		CreatedAt:    "2009-05-27 16:29:46",
		PushedAt:     "2014-12-31 21:16:33",
	}
	if got := s.RepoByShortPath(want.ShortPath); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestRecs(t *testing.T) {
	s := NewInMemRepoStore()

	if err := s.LoadRecsFromCSVGZip(recsCSVGZip); err != nil {
		t.Fatalf("LoadRecsFromCSVGZip(%q) failed: %v", recsCSVGZip, err)
	}
	want := []RepoRelation{
		{"google/material-design-icons", 820.0},
		{"audreyr/favicon-cheat-sheet", 313.0},
		{"IanLunn/jQuery-Parallax", 56},
	}
	if got := s.SimilarRepos("AFNetworking/AFNetworking"); !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
