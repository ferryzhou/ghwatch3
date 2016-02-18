package models

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Repo struct {
	ShortPath    string
	URL          string
	Stars        int
	Language     string
	Description  string
	Website      string
	name         string
	owner        string
	organization string
	createdAt    string
	pushedAt     string
}

type RepoStore interface {
	SimilarRepos(repo *Repo) []*Repo
	RepoByShortPath(s string) *Repo
}

type InMemRepoStore struct {
	repos map[string]*Repo
	rec   map[string][]string
}

func (r *InMemRepoStore) SimilarRepos(repo *Repo) []*Repo {
	ps := r.rec[repo.ShortPath]
	rs := make([]*Repo, 0, len(ps))
	for _, p := range ps {
		rs = append(rs, r.repos[p])
	}
	return rs
}

func (r *InMemRepoStore) RepoByShortPath(s string) *Repo {
	return r.repos[s]
}

func (r *InMemRepoStore) ReposShortPaths(ss []string) []*Repo {
	repos := []*Repo{}
	for _, s := range ss {
		if r.repos[s] != nil {
			repos = append(repos, r.repos[s])
		}
	}
	return repos
}

func NewInMemRepoStore() *InMemRepoStore {
	return &InMemRepoStore{make(map[string]*Repo), make(map[string][]string)}
}

// InMemRepoStoreFromCSV creates a repo store from raw csv files
func InMemRepoStoreFromCSV(reposPath, recPath string) (*InMemRepoStore, error) {
	s := NewInMemRepoStore()
	f, err := os.Open(reposPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %v: %v", reposPath, err)
	}
	defer f.Close()
	r := csv.NewReader(bufio.NewReader(f))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read file %v: %v", reposPath, err)
		}
		repo, err := repoFromCSVRecord(record)
		if err != nil {
			log.Printf("Not valid repo record %v: %v", record, err)
			continue
		}
		s.repos[repo.ShortPath] = repo
	}
	return s, nil
}

func repoFromCSVRecord(r []string) (*Repo, error) {
	if len(r) < 9 {
		return nil, fmt.Errorf("not enough length, require at least 9, got %v", len(r))
	}

	name := r[0]
	owner := r[1]
	org := r[2]
	lang := r[3]
	url := r[4]
	created := r[5]
	desc := r[6]
	watchers, _ := strconv.Atoi(r[7])
	pushed := r[8]

	if len(url) < prefixLen {
		return nil, fmt.Errorf("url %q not valid", url)
	}
	return &Repo{
		ShortPath:    ShortPathFromURL(url),
		URL:          url,
		Stars:        watchers,
		Language:     lang,
		Description:  desc,
		name:         name,
		owner:        owner,
		organization: org,
		createdAt:    created,
		pushedAt:     pushed,
	}, nil
}

var prefixLen = len("https://github.com/")

func ShortPathFromURL(url string) string {
	return url[prefixLen:]
}

func (r *InMemRepoStore) Save(path string) error {
	return nil
}
