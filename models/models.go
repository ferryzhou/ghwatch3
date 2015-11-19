package models

type Repo struct {
	Id          uint64
	ShortPath   string
	URL         string
	Stars       int
	Language    string
	Description string
	Website     string
}

type RepoIndex struct {
	sp2id map[string]uint64
	id2sp map[uint64]string
}

type RepoStore struct {
	index *RepoIndex
	repos map[uint64]*Repo
	// Recommendations, id to ids mapping
	rec map[uint64][]uint64
}

func (r *RepoStore) GetSimilarRepos(repo *Repo) []*Repo {

}

func (r *RepoStore) GetRepoByIds(ids []uint64) []*Repo {
}

func InitRepoStore() *RepoStore {
}

func RepoStoreFromCSV(path string) (*RepoStore, error) {
}

func (r *RepoStore) Save(path string) error {
}
