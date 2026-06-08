package main

import "testing"

func TestReadDir(t *testing.T) {
	// Place your code here
	_, err := ReadDir("testdata/env/")
	if err != nil {
		t.Fatal(err)
	}

}
