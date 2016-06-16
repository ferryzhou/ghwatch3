package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ferryzhou/ghwatch3/models"
	"github.com/gorilla/mux"
)

var (
	s               models.RepoStore
	reposPath       = flag.String("repos_path", "", "repos csv gzip file path")
	recsPathPattern = flag.String("recs_path_pattern", "", "recs csv gzip path pattern")
	fromGobPath     = flag.String("from_gob_path", "", "in mem store gob data as input, if specified, load data from here instead from csv")
	toGobPath       = flag.String("to_gob_path", "", "in mem store gob data as output")
)

func initStore() {
	ms := &models.InMemRepoStore{}
	if *fromGobPath != "" {
		if err := ms.Load(*fromGobPath); err != nil {
			log.Fatalf("failed to load: %v", err)
		}
	} else {
		ms := models.NewInMemRepoStore()
		if err := ms.LoadReposFromCSVGZip(*reposPath); err != nil {
			log.Fatalf("LoadReposFromCSVGZip(%q) failed: %v", *reposPath, err)
		}
		if err := ms.LoadRecsFromCSVGZipFiles(*recsPathPattern); err != nil {
			log.Fatalf("LoadRecsFromCSVGZip(%q) failed: %v", *recsPathPattern, err)
		}
		ms.GenReposByOwner()
	}

	s = ms
	log.Printf("Repo store is ready")

	if *toGobPath != "" {
		log.Printf("Writing to gob file %v", *toGobPath)
		if err := ms.Save(*toGobPath); err != nil {
			log.Fatalf("Failed to save: %v", err)
		}
	}
}
func main() {
	flag.Parse()
	initStore()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/repos", RepoIndex)
	router.HandleFunc("/repos/{owner}", OwnerShow)
	router.HandleFunc("/repos/{owner}/{repo}", RepoShow)
	router.HandleFunc("/recs/{owner}/{repo}", RepoRec)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func RepoIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Repo Index!")
}

func RepoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(w, "Repo show:", vars["owner"], "/", vars["repo"])
	sp := vars["owner"] + "/" + vars["repo"]
	rp := s.RepoByShortPath(sp)
	if err := json.NewEncoder(w).Encode(rp); err != nil {
		log.Printf("encode error %v: %v", rp, err)
	}
}

func RepoRec(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sp := vars["owner"] + "/" + vars["repo"]
	rp := s.SimilarRepos(sp)
	if err := json.NewEncoder(w).Encode(rp); err != nil {
		log.Printf("encode error %v: %v", rp, err)
	}
}

func OwnerShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rp := s.ReposByOwner(vars["owner"])
	if err := json.NewEncoder(w).Encode(rp); err != nil {
		log.Printf("encode error %v: %v", rp, err)
	}
}
