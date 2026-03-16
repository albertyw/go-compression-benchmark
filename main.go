package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type dataFile struct {
	name string
	path string
}

func main() {
	files := []dataFile{
		{name: "data.h3", path: "data.h3"},
		{name: "data_mock.h3", path: "data_mock.h3"},
	}

	var sections []ResultSection
	for _, f := range files {
		data, err := os.ReadFile(filepath.Clean(f.path))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: cannot read %s: %v\n", f.path, err)
			os.Exit(1)
		}
		fmt.Printf("Benchmarking %s (%d bytes)...\n", f.name, len(data))
		results := RunBenchmarks(data)
		sections = append(sections, ResultSection{
			FileName: f.name,
			FileSize: len(data),
			Results:  results,
		})
	}

	if err := WriteResults(sections); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cannot write RESULTS.md: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Results written to RESULTS.md")
}
