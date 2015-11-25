package models

type Repo struct {
	ShortPath   string
	URL         string
	Stars       int
	Language    string
	Description string
	Website     string
}

type RepoStore struct {
	repos map[string]*Repo
	rec map[string][]string
}

func (r *RepoStore) GetSimilarRepos(repo *Repo) []*Repo {
  
}

func (r *RepoStore) GetRepoByShortPath(s string) *Repo {
	return r.repos[s]
}

func (r *RepoStore) GetRepoByShortPaths(ss []string) []*Repo {
	repos := []*Repo{}
	for _, s := range ss {
	  if r.repos[s] != nil {
			repos = append(repos, r.repos[s])
		}
	}
	return repos
}

func InitRepoStore() *RepoStore {
}
// RepoStoreFromCSV creates a repo store from raw csv files
// 1) repos.csv, 
// 2) rec.csv
func RepoStoreFromCSV(path string) (*RepoStore, error) {
}

func (r *RepoStore) Save(path string) error {
}
