package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ferryzhou/ghwatch3/models"
	"github.com/gorilla/mux"
)

var s models.RepoStore

func main() {
	repoCSV := "../models/testdata/repos.csv"
	recCSV := ""
	ss, err := models.InMemRepoStoreFromCSV(repoCSV, recCSV)
	if err != nil {
		log.Fatalf("InMemRepoStoreFromCSV(%q, %q) failed: %v", repoCSV, recCSV, err)
	}
	s = ss
	log.Printf("Repo store is ready")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/repos", RepoIndex)
	router.HandleFunc("/repos/{owner}", OwnerShow)
	router.HandleFunc("/repos/{owner}/{repo}", RepoShow)
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

func OwnerShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(w, "Owner show:", vars["owner"])
}
