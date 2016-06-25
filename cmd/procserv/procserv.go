// Program procserv is a http server provides api for processed data.
// Example:
//   cd ghwatch3
//   go run cmd/procserv/procserv.go --in_gob_path=processed/recs.gob
//   curl http://localhost:8080/rec/twbs/bootstrap
//   curl http://localhost:8080/recn/twbs/bootstrap
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ferryzhou/ghwatch3/proc"
	"github.com/gorilla/mux"
)

var (
	p         *proc.Proc
	inGobPath = flag.String("in_gob_path", "", "in mem store gob data as input")
	port      = flag.String("port", "8080", "port number")
)

func main() {
	flag.Parse()
	p = &proc.Proc{}
	if err := p.Load(*inGobPath); err != nil {
		log.Fatalf("failed to load: %v", err)
	}
	log.Printf("Data loaded (%v), ready to serve ...", p)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/repos", repoIndex)
	router.HandleFunc("/rec/{owner}/{repo}", repoRecRaw)
	router.HandleFunc("/recn/{owner}/{repo}", repoRecNorm)
	log.Fatal(http.ListenAndServe(":"+*port, router))
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func repoIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Repo Index!")
}

func repoRecNorm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sp := vars["owner"] + "/" + vars["repo"]
	rp := p.Rec(sp)
	if err := json.NewEncoder(w).Encode(rp); err != nil {
		log.Printf("encode error %v: %v", rp, err)
	}
}

func repoRecRaw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sp := vars["owner"] + "/" + vars["repo"]
	rp := p.RecRaw(sp)
	bs, err := json.Marshal(rp)
	if err != nil {
		log.Printf("encode error %v: %v", rp, err)
	}
	callback := r.FormValue("callback")
	if callback != "" {
		fmt.Fprintf(w, "%s(%s)", callback, bs)
	} else {
		w.Write(bs)
	}
}
