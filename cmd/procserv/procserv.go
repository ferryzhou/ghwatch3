// Program procserv is a http server provides api for processed data.
// Example:
//   cd ghwatch3
//   go run cmd/procserv/procserv.go --in_gob_path=processed/recs.gob
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ferryzhou/ghwatch3/proc"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
	router.HandleFunc("/repos", repoIndex)
	router.HandleFunc("/rec", repoRec)
	router.HandleFunc("/hello", hello)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
	http.Handle("/", router)
	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST).
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":"+*port, handler))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func repoIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Repo Index!")
}

// /rec?sp=<short_path>
// /rec?sp=<short_path>&norm=true
func repoRec(w http.ResponseWriter, r *http.Request) {
	sp := r.FormValue("sp")
	rp := p.RecRaw(sp)
	log.Printf("norm: %q", r.FormValue("norm"))
	if r.FormValue("norm") != "" {
		rp = p.RecNorm(sp)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rp); err != nil {
		log.Printf("encode error %v: %v", rp, err)
	}
}
