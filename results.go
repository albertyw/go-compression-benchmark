package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/albertyw/go-compression-benchmark/compressors"
)

// ResultSection groups benchmark results for one data file.
type ResultSection struct {
	FileName string
	FileSize int
	Results  []BenchmarkResult
}

// WriteResults writes benchmark results for one data file to a results file.
// Returns the path of the written file.
func WriteResults(section ResultSection) (string, error) {
	var sb strings.Builder
	sb.WriteString("# Compression Benchmark Results\n\n")

	// Environment info
	sb.WriteString("## Environment\n\n")
	sb.WriteString(fmt.Sprintf("- **CPU:** %s\n", getCPUModel()))
	sb.WriteString(fmt.Sprintf("- **Architecture:** %s\n", runtime.GOARCH))
	sb.WriteString(fmt.Sprintf("- **OS:** %s\n", runtime.GOOS))
	sb.WriteString(fmt.Sprintf("- **Go:** %s\n", runtime.Version()))
	sb.WriteString("\n")

	sb.WriteString(fmt.Sprintf("## %s (%s bytes)\n\n", section.FileName, formatInt(section.FileSize)))

	// Compression table
	sb.WriteString("### Compression\n\n")
	sb.WriteString("| Library | Level | Compressed Size | Ratio | Time | Memory | Allocs | Pure Go | Notes |\n")
	sb.WriteString("|---------|-------|----------------|-------|------|--------|--------|---------|-------|\n")
	for _, r := range section.Results {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
			libraryLink(r.Compressor.Name, r.Compressor.Library),
			r.Compressor.Level,
			formatBytes(r.CompressedSize),
			formatRatio(r.Ratio),
			formatDuration(r.CompressTime),
			formatBytes(int(r.CompressMem)),
			formatInt(int(r.CompressAllocs)),
			formatPureGo(r.Compressor),
			r.Compressor.Notes,
		))
	}
	sb.WriteString("\n")

	// Decompression table
	sb.WriteString("### Decompression\n\n")
	sb.WriteString("| Library | Level | Time | Memory | Allocs |\n")
	sb.WriteString("|---------|-------|------|--------|--------|\n")
	for _, r := range section.Results {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			libraryLink(r.Compressor.Name, r.Compressor.Library),
			r.Compressor.Level,
			formatDuration(r.DecompressTime),
			formatBytes(int(r.DecompressMem)),
			formatInt(int(r.DecompressAllocs)),
		))
	}
	sb.WriteString("\n")

	sb.WriteString("## Notes\n\n")
	sb.WriteString("- All libraries are pure Go with no CGo dependencies\n")
	sb.WriteString("- Benchmarks run 2 warmup + 10 measured iterations per compressor; median values reported\n")
	sb.WriteString("- Times measured via monotonic clock (time.Since) in timing-only runs\n")
	sb.WriteString("- Memory/allocs measured via runtime.MemStats deltas in separate memory-only runs (after runtime.GC())\n")
	sb.WriteString("- Rows sorted by compressed size ascending within each table\n")
	sb.WriteString("- Ratio = original size / compressed size (higher is better)\n")
	sb.WriteString("- Memory figures are approximate; runtime.MemStats uses process-global counters, so allocations from concurrent goroutines may be included\n")

	path := resultsPath(section.FileName)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	return path, os.WriteFile(path, []byte(sb.String()), 0644)
}

// resultsPath returns the output file path, e.g. "results/data-2026-03-15-amd64-linux.md".
func resultsPath(dataFileName string) string {
	base := strings.TrimSuffix(dataFileName, filepath.Ext(dataFileName))
	date := time.Now().Format("2006-01-02")
	return filepath.Join("results", fmt.Sprintf("%s-%s-%s-%s.md", base, date, runtime.GOARCH, runtime.GOOS))
}

func formatPureGo(c compressors.Compressor) string {
	if c.Library == "" {
		return "—"
	}
	if c.PureGo {
		return "Yes"
	}
	return "No"
}

// libraryLink returns a markdown link for the library name pointing to GitHub.
func libraryLink(name, library string) string {
	url := libraryURL(library)
	if url == "" {
		return name
	}
	return fmt.Sprintf("[%s](%s)", name, url)
}

// libraryURL derives a GitHub or pkg.go.dev URL from a Go import path.
func libraryURL(library string) string {
	if library == "" {
		return ""
	}
	// Standard library packages
	if !strings.Contains(library, ".") {
		return "https://pkg.go.dev/" + library
	}
	// GitHub packages: extract github.com/owner/repo
	if strings.HasPrefix(library, "github.com/") {
		parts := strings.Split(library, "/")
		if len(parts) >= 3 {
			repo := strings.Join(parts[:3], "/")
			// Strip version suffixes like /v4
			return "https://" + repo
		}
	}
	return ""
}

// getCPUModel returns the CPU model name.
func getCPUModel() string {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "unknown"
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "unknown"
}

func formatBytes(b int) string {
	switch {
	case b >= 1024*1024:
		return fmt.Sprintf("%.2f MB", float64(b)/(1024*1024))
	case b >= 1024:
		return fmt.Sprintf("%.2f KB", float64(b)/1024)
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func formatInt(n int) string {
	if n < 0 {
		return "-" + formatInt(-n)
	}
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	s := fmt.Sprintf("%d", n)
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}

func formatRatio(r float64) string {
	return fmt.Sprintf("%.2fx", r)
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "—"
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.1fµs", float64(d)/float64(time.Microsecond))
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	}
	return fmt.Sprintf("%.2fs", float64(d)/float64(time.Second))
}
