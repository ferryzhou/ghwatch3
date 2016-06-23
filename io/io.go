// Package io provides functions for data input/output.
package io

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
)

// SaveGob saves data to a file with gob encoder.
func SaveGob(path string, d interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %v: %v", path, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	enc := gob.NewEncoder(w)
	if err = enc.Encode(d); err != nil {
		return fmt.Errorf("failed to encode data: %v", err)
	}
	return nil
}

// LoadGob load data from a file with gob decoder.
func LoadGob(path string, d interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %v: %v", path, err)
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	if err = dec.Decode(d); err != nil {
		return fmt.Errorf("failed to decode: %v", err)
	}
	return nil
}
