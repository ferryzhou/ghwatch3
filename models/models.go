package models

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
	CreatedAt    string
	PushedAt     string
}

type RepoStore interface {
	SimilarRepos(shortPath string) []RepoRelation
	RepoByShortPath(s string) *Repo
	ReposByOwner(s string) []*Repo
}

type InMemRepoStore struct {
	Repos      map[string]*Repo
	Recs       map[string][]RepoRelation
	OwnerRepos map[string][]*Repo
}

type RepoRelation struct {
	ShortPath string
	Score     float64
}

type ByScore []RepoRelation

func (a ByScore) Len() int      { return len(a) }
func (a ByScore) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// sort so that higher score is ranked higher.
func (a ByScore) Less(i, j int) bool { return a[i].Score >= a[j].Score }

func (r *InMemRepoStore) SimilarRepos(shortPath string) []RepoRelation {
	return r.Recs[shortPath]
}

func (r *InMemRepoStore) RepoByShortPath(s string) *Repo {
	return r.Repos[s]
}

func (r *InMemRepoStore) ReposShortPaths(ss []string) []*Repo {
	repos := []*Repo{}
	for _, s := range ss {
		if r.Repos[s] != nil {
			repos = append(repos, r.Repos[s])
		}
	}
	return repos
}

func (r *InMemRepoStore) ReposByOwner(s string) []*Repo {
	return r.OwnerRepos[s]
}

func NewInMemRepoStore() *InMemRepoStore {
	return &InMemRepoStore{make(map[string]*Repo), make(map[string][]RepoRelation), make(map[string][]*Repo)}
}

func (r *InMemRepoStore) GenReposByOwner() {
	for _, v := range r.Repos {
		r.OwnerRepos[v.owner] = append(r.OwnerRepos[v.owner], v)
	}
}

func ReposFromCSVReader(r *csv.Reader, repos map[string]*Repo) error {
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read csv record: %v", err)
		}
		repo, err := repoFromCSVRecord(record)
		if err != nil {
			log.Printf("not valid repo record %v: %v", record, err)
			continue
		}
		sr, ok := repos[repo.ShortPath]
		if !ok || strings.Compare(sr.PushedAt, repo.PushedAt) < 0 {
			repos[repo.ShortPath] = repo
		}
	}
	log.Printf("%v repos has been loaded", len(repos))
	return nil
}

func ReposFromCSVGZip(path string, repos map[string]*Repo) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %v: %v", path, err)
	}
	defer f.Close()
	ar, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to unzip %v: %v", path, err)
	}
	defer ar.Close()

	r := csv.NewReader(bufio.NewReader(ar))
	return ReposFromCSVReader(r, repos)
}

func ReposFromCSVFile(path string, repos map[string]*Repo) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %v: %v", path, err)
	}
	defer f.Close()
	r := csv.NewReader(bufio.NewReader(f))
	return ReposFromCSVReader(r, repos)
}

func (s *InMemRepoStore) LoadReposFromCSVGZip(path string) error {
	return ReposFromCSVGZip(path, s.Repos)
}

func (s *InMemRepoStore) LoadRecsFromCSVGZip(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %v: %v", path, err)
	}
	defer f.Close()
	ar, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to unzip %v: %v", path, err)
	}
	defer ar.Close()

	r := csv.NewReader(bufio.NewReader(ar))
	return RecsFromCSVReader(r, s.Recs)
}
func (s *InMemRepoStore) LoadRecsFromCSVGZipFiles(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to read dir %v: %v", pattern, err)
	}
	for _, file := range files {
		if err = s.LoadRecsFromCSVGZip(file); err != nil {
			return err
		}
	}
	return nil
}

// InMemRepoStoreFromCSV creates a repo store from raw csv files
func InMemRepoStoreFromCSV(reposPath, recPath string) (*InMemRepoStore, error) {
	s := NewInMemRepoStore()
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
		CreatedAt:    created,
		PushedAt:     pushed,
	}, nil
}

func RecsFromCSVReader(r *csv.Reader, recs map[string][]RepoRelation) error {
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read csv record: %v", err)
		}
		if len(record[0]) < prefixLen {
			log.Printf("not valid rec record %v", record)
			continue
		}
		sp1 := ShortPathFromURL(record[0])
		sp2 := ShortPathFromURL(record[1])
		c, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Printf("Failed to ParseFloat(%q, 64): %v", record[2], err)
		}
		recs[sp1] = append(recs[sp1], RepoRelation{sp2, c})
	}
	for k, _ := range recs {
		sort.Sort(ByScore(recs[k]))
	}
	log.Printf("%v recs has been loaded", len(recs))
	return nil
}

var prefixLen = len("https://github.com/")

func ShortPathFromURL(url string) string {
	return url[prefixLen:]
}

func (r *InMemRepoStore) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %v: %v", path, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	enc := gob.NewEncoder(w)
	if err = enc.Encode(r); err != nil {
		return fmt.Errorf("failed to encode data: %v", err)
	}
	return nil
}

func (r *InMemRepoStore) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %v: %v", path, err)
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	if err = dec.Decode(r); err != nil {
		return fmt.Errorf("failed to decode: %v", err)
	}
	return nil
}
