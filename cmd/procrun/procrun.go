package main

import (
	"flag"
	"log"
	"os"

	"github.com/ferryzhou/ghwatch3/proc"
)

var (
	inDir  = flag.String("in_dir", "", "input directory")
	outDir = flag.String("out_dir", "", "output directory")
)

// cd ghwatch3
// go run cmd/procrun/procrun.go --in_dir=results --out_dir=processed
func main() {
	flag.Parse()
	p := proc.NewProc(*inDir, *outDir, "c1000")
	if err := os.MkdirAll(*outDir, 0777); err != nil {
		log.Fatalf("Failed to create folder %v: %v", *outDir, err)
	}
	if err := p.Run(); err != nil {
		log.Fatalf("Failed to process data: %v", err)
	}
	if err := p.Save(); err != nil {
		log.Fatalf("Failed to save results: %v", err)
	}
	log.Printf("Successfully processed data .....")
}
