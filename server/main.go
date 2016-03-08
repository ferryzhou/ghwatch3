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
)

func main() {
	flag.Parse()
	ss := models.NewInMemRepoStore()
	if err := ss.LoadReposFromCSVGZip(*reposPath); err != nil {
		log.Fatalf("LoadReposFromCSVGZip(%q) failed: %v", *reposPath, err)
	}
	if err := ss.LoadRecsFromCSVGZipFiles(*recsPathPattern); err != nil {
		log.Fatalf("LoadRecsFromCSVGZip(%q) failed: %v", *recsPathPattern, err)
	}
	s = ss
	log.Printf("Repo store is ready")
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
	fmt.Fprintln(w, "Owner show:", vars["owner"])
}
