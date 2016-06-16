package main

import (
	"compress/gzip"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	inDir  = flag.String("in_dir", "", "input directory")
	outDir = flag.String("out_dir", "", "output directory")
)

// Proc provides functions to process the raw bigquery data.
type Proc struct {
	shortPathes                 []string
	shortPathToIDMap            map[string]int
	inDir, outDir               string
	rawReposPath, procReposPath string
}

// NewProc creates a new Proc object given inDir, outDir and n, e.g. c1000.
// rawReposPath is <inDir>/repos_<n>_full.csv.gz
func NewProc(inDir, outDir, n string) *Proc {
	return &Proc{
		inDir:            inDir,
		outDir:           outDir,
		rawReposPath:     filepath.Join(inDir, "repos_"+n+"_full.csv.gz"),
		procReposPath:    filepath.Join(outDir, "repos.csv"),
		shortPathToIDMap: make(map[string]int),
	}
}

func (p *Proc) procRepos() error {
	f, err := os.Open(p.rawReposPath)
	if err != nil {
		return fmt.Errorf("failed to open file %v: %v", p.rawReposPath, err)
	}
	defer f.Close()
	ar, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to unzip %v: %v", p.rawReposPath, err)
	}
	defer ar.Close()
	r := csv.NewReader(ar)

	of, err := os.Create(p.procReposPath)
	if err != nil {
		return fmt.Errorf("failed to create file %v: %v", p.procReposPath, err)
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
		if _, ok := p.shortPathToIDMap[shortPath]; ok {
			continue
		}
		id := len(p.shortPathes)
		p.shortPathToIDMap[shortPath] = id
		or := []string{strconv.Itoa(id), shortPath}
		or = append(or, record...)
		if err := w.Write(or); err != nil {
			return fmt.Errorf("failed to write record: %v", err)
		}
		p.shortPathes = append(p.shortPathes, shortPath)
	}
	log.Printf("%v repos has been loaded", len(p.shortPathes))

	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("failed flush file: %v", err)
	}
	return nil
}

var prefixLen = len("https://github.com/")

func shortPathFromURL(url string) string {
	return url[prefixLen:]
}

func main() {
	flag.Parse()
	p := NewProc(*inDir, *outDir, "c1000")
	if err := os.MkdirAll(*outDir, 0777); err != nil {
		log.Fatalf("Failed to create folder %v: %v", *outDir, err)
	}
	if err := p.procRepos(); err != nil {
		log.Fatalf("Failed to process data: %v", err)
	}
}
