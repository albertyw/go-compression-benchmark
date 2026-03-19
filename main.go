package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const sampleDataDir = "sample_data"

type dataFile struct {
	name string
	path string
}

func discoverDataFiles() ([]dataFile, error) {
	entries, err := os.ReadDir(sampleDataDir)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", sampleDataDir, err)
	}
	var files []dataFile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		files = append(files, dataFile{
			name: e.Name(),
			path: filepath.Join(sampleDataDir, e.Name()),
		})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].name < files[j].name })
	return files, nil
}

func main() {
	files, err := discoverDataFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: no data files found in %s\n", sampleDataDir)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "ERROR: cannot write results: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Results written to %s\n", resultsPath())
}
