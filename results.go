package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ResultSection groups benchmark results for one data file.
type ResultSection struct {
	FileName string
	FileSize int
	Results  []BenchmarkResult
}

// WriteResults writes benchmark results to RESULTS.md.
func WriteResults(sections []ResultSection) error {
	var sb strings.Builder
	sb.WriteString("# Compression Benchmark Results\n\n")

	for _, section := range sections {
		sb.WriteString(fmt.Sprintf("## %s (%s bytes)\n\n", section.FileName, formatInt(section.FileSize)))
		sb.WriteString("| Library | Level | Compressed Size | Ratio | Compress Time | Decompress Time | Compress Mem | Compress Allocs | Decompress Mem | Decompress Allocs | Pure Go | Notes |\n")
		sb.WriteString("|---------|-------|----------------|-------|--------------|----------------|-------------|----------------|---------------|-------------------|---------|-------|\n")

		for _, r := range section.Results {
			pureGo := "—"
			if r.Compressor.Library != "" {
				if r.Compressor.PureGo {
					pureGo = "Yes"
				} else {
					pureGo = "No"
				}
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s |\n",
				r.Compressor.Name,
				r.Compressor.Level,
				formatBytes(r.CompressedSize),
				formatRatio(r.Ratio),
				formatDuration(r.CompressTime),
				formatDuration(r.DecompressTime),
				formatBytes(int(r.CompressMem)),
				formatInt(int(r.CompressAllocs)),
				formatBytes(int(r.DecompressMem)),
				formatInt(int(r.DecompressAllocs)),
				pureGo,
				r.Compressor.Notes,
			))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Notes\n\n")
	sb.WriteString("- All libraries are pure Go with no CGo dependencies\n")
	sb.WriteString("- Benchmarks run 2 warmup + 10 measured iterations per compressor; median values reported\n")
	sb.WriteString("- Times measured via monotonic clock (time.Since) in timing-only runs\n")
	sb.WriteString("- Memory/allocs measured via runtime.MemStats deltas in separate memory-only runs (after runtime.GC())\n")
	sb.WriteString("- Rows sorted by compressed size ascending within each table\n")
	sb.WriteString("- Ratio = original size / compressed size (higher is better)\n")
	sb.WriteString("- Memory figures are approximate; runtime.MemStats uses process-global counters, so allocations from concurrent goroutines may be included\n")

	return os.WriteFile("RESULTS.md", []byte(sb.String()), 0644)
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
