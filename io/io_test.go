package io

import (
	"reflect"
	"strconv"
	"testing"
)

func TestSaveLoad(t *testing.T) {
	testCases := []struct {
		name  string
		value string
	}{
		{
			name:  "string",
			value: "abc",
		},
	}
	for i, tc := range testCases {
		path := "/tmp/ghwatch3_io_test_" + strconv.Itoa(i) + ".gob"
		if err := SaveGob(path, tc.value); err != nil {
			t.Fatalf("%s: failed to save data %v: %v", tc.name, tc.value, err)
		}
		got := ""
		if err := LoadGob(path, &got); err != nil {
			t.Fatalf("%s: failed to load data from %v: %v", tc.name, path, err)
		}
		if !reflect.DeepEqual(got, tc.value) {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.value)
		}
	}
}
