package proc

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	gio "github.com/ferryzhou/ghwatch3/io"
)

// Proc provides functions to process the raw bigquery data.
type Proc struct {
	// Keep top N related repos
	TopN int
	// short path of i-th repo
	ShortPathes []string
	// map from short path to index
	ShortPathToIDMap map[string]int
	// recommends for i-th repo based on raw overlap count
	RecsRaw [][]int
	// recommends for i-th repo based on normalized count
	RecsNorm [][]int
	// input and output directory
	inDir, outDir                  string
	repoInPath, repoOutCSVPath     string
	recsInPattern, modelOutGobPath string
}

// NewProc creates a new Proc object given inDir, outDir and n, e.g. c1000.
// repoInPath is <inDir>/repos_<n>_full.csv.gz
func NewProc(inDir, outDir, n string) *Proc {
	return &Proc{
		TopN:             20,
		inDir:            inDir,
		outDir:           outDir,
		repoInPath:       filepath.Join(inDir, "repos_"+n+"_full.csv.gz"),
		repoOutCSVPath:   filepath.Join(outDir, "repos.csv"),
		recsInPattern:    filepath.Join(inDir, "repo_repo_count_"+n+"_*.csv.gz"),
		modelOutGobPath:  filepath.Join(outDir, "recs.gob"),
		ShortPathToIDMap: make(map[string]int),
	}
}

// Run runs the processing pipeline.
func (p *Proc) Run() error {
	if err := p.procRepos(); err != nil {
		return fmt.Errorf("Failed to process data: %v", err)
	}
	if err := p.procRecs(); err != nil {
		return fmt.Errorf("Failed to process recs data :%v", err)
	}
	return nil
}

// Save saves data to a file.
func (p *Proc) Save() error {
	return gio.SaveGob(p.modelOutGobPath, p)
}

// Load loads data from a file.
func (p *Proc) Load(path string) error {
	return gio.LoadGob(path, p)
}

// Rec return recommended repos given a repo's short path.
func (p *Proc) Rec(shortPath string) []string {
	return p.recWith(shortPath, p.RecsNorm)
}

// RecRaw return recommended repos given a repo's short path.
func (p *Proc) RecRaw(shortPath string) []string {
	return p.recWith(shortPath, p.RecsRaw)
}

func (p *Proc) recWith(shortPath string, recs [][]int) []string {
	ind, ok := p.ShortPathToIDMap[shortPath]
	if !ok {
		return nil
	}
	log.Printf("input: %v, index: %v", shortPath, ind)
	rs := make([]string, 0, len(recs[ind]))
	for _, r := range recs[ind] {
		rs = append(rs, p.ShortPathes[r])
	}
	return rs
}

func (p *Proc) String() string {
	return fmt.Sprintf("%v repos, %v recs", len(p.ShortPathes), len(p.RecsNorm))
}

// procRepos process raw repos data, including
//   1. remove duplicated records
//   2. add id (based on popularity)
//   3. add shortPath, e.g. twitter/bootstrap
// The result is a csv file without header.
func (p *Proc) procRepos() error {
	f, err := os.Open(p.repoInPath)
	if err != nil {
		return fmt.Errorf("failed to open file %v: %v", p.repoInPath, err)
	}
	defer f.Close()
	ar, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to unzip %v: %v", p.repoInPath, err)
	}
	defer ar.Close()
	r := csv.NewReader(ar)

	of, err := os.Create(p.repoOutCSVPath)
	if err != nil {
		return fmt.Errorf("failed to create file %v: %v", p.repoOutCSVPath, err)
	}
	w := csv.NewWriter(of)

	// skip header
	r.Read()
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read csv record: %v", err)
		}
		url := record[4]
		shortPath := shortPathFromURL(url)
		if _, ok := p.ShortPathToIDMap[shortPath]; ok {
			continue
		}
		id := len(p.ShortPathes)
		p.ShortPathToIDMap[shortPath] = id
		or := []string{strconv.Itoa(id), shortPath}
		or = append(or, record...)
		if err := w.Write(or); err != nil {
			return fmt.Errorf("failed to write record: %v", err)
		}
		p.ShortPathes = append(p.ShortPathes, shortPath)
	}
	log.Printf("%v repos has been loaded", len(p.ShortPathes))

	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("failed flush file: %v", err)
	}
	return nil
}

// procRecs process raw repo relations data:
//   1. replace url with id
//   2. for each repo, find top N related repos.
// Two scoring methods:
//  1. raw overlap count, say Co.
//  2. normalized value, i.e. Co / (sqrt(C1) * sqrt(C2))
func (p *Proc) procRecs() error {
	return p.procRecsFromFilePattern(p.recsInPattern)
}

func (p *Proc) procRecsFromFilePattern(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to read dir %v: %v", pattern, err)
	}
	for _, file := range files {
		if err = p.procRecsFromCSVGZip(file); err != nil {
			return err
		}
	}
	return nil
}

func (p *Proc) procRecsFromCSVGZip(path string) error {
	log.Printf("Process file %q", path)
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

	r := csv.NewReader(ar)
	return p.procRecsFromCSVReader(r)
}

//
type repoRelation struct {
	i     int
	score float64
}

type byScore []repoRelation

func (a byScore) Len() int      { return len(a) }
func (a byScore) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// sort so that higher score is ranked higher.
func (a byScore) Less(i, j int) bool { return a[i].score >= a[j].score }

func (p *Proc) repoIndexFromURL(url string) int {
	return p.ShortPathToIDMap[shortPathFromURL(url)]
}

func (p *Proc) procRecsFromCSVReader(r *csv.Reader) error {
	repoCount := len(p.ShortPathes)
	relationRaw := make([][]repoRelation, repoCount)
	relationNorm := make([][]repoRelation, repoCount)
	selfCounts := make([]float64, repoCount)
	n := 0
	r.Read()
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read csv record: %v", err)
		}
		repo1 := p.repoIndexFromURL(record[0])
		repo2 := p.repoIndexFromURL(record[1])
		count, err := strconv.Atoi(record[2])
		if err != nil {
			log.Printf("Failed to ParseFloat(%q, 64): %v", record[2], err)
		}
		relationRaw[repo1] = append(relationRaw[repo1], repoRelation{repo2, float64(count)})
		if repo1 == repo2 {
			selfCounts[repo1] = float64(count)
		}
		n++
		if math.Mod(float64(n), float64(500000)) == 0 {
			log.Printf("Processed %v recs", n)
		}
	}

	for i, rs := range relationRaw {
		for _, r := range rs {
			normScore := r.score / math.Sqrt(selfCounts[i]) / math.Sqrt(selfCounts[r.i])
			relationNorm[i] = append(relationNorm[i], repoRelation{r.i, normScore})
		}
	}

	p.RecsRaw = repoRelationsToRecs(relationRaw, p.TopN)
	p.RecsNorm = repoRelationsToRecs(relationNorm, p.TopN)

	return nil
}

func repoRelationsToRecs(rs [][]repoRelation, topN int) [][]int {
	out := make([][]int, 0, len(rs))
	for _, rr := range rs {
		sort.Sort(byScore(rr))
		sr := []int{}
		for j, r := range rr {
			if j == topN {
				break
			}
			sr = append(sr, r.i)
		}
		out = append(out, sr)
	}
	return out
}

var prefixLen = len("https://github.com/")

func shortPathFromURL(url string) string {
	return url[prefixLen:]
}
