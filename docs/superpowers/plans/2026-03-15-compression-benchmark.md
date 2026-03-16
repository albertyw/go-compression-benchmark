# Compression Benchmark Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go program that benchmarks 11 compression libraries against H3 timezone data and outputs a markdown results table.

**Architecture:** Registry pattern with Compressor structs. Each library registers via init(). Benchmark runner measures timing and memory in separate passes, takes medians, and writes RESULTS.md. All files in `package main` (root directory).

**Tech Stack:** Go, klauspost/compress (gzip/deflate/zstd/snappy/s2/pgzip), andybalholm/brotli, pierrec/lz4/v4, ulikunitz/xz, golang/snappy

**Spec:** `docs/superpowers/specs/2026-03-15-compression-benchmark-design.md`

---

## File Structure

All `.go` files are in the root directory, `package main`.

| File | Responsibility |
|------|---------------|
| `go.mod` | Module definition and dependencies |
| `compressor.go` | Compressor struct (with Name + Level fields) + registry (Register/GetAll) |
| `benchmark.go` | Benchmark runner: warmup, timing runs, memory runs, median calculation |
| `results.go` | Markdown table formatter, writes RESULTS.md |
| `main.go` | Entry point: load data, run benchmarks, write results |
| `compressor_uncompressed.go` | Baseline (no-op) compressor |
| `compressor_gzip_stdlib.go` | stdlib gzip at 3 levels |
| `compressor_gzip_klauspost.go` | klauspost gzip at 3 levels |
| `compressor_deflate_klauspost.go` | klauspost deflate at 3 levels |
| `compressor_zstd.go` | zstd at 4 levels |
| `compressor_brotli.go` | brotli at 3 levels |
| `compressor_snappy.go` | golang/snappy (single level) |
| `compressor_s2.go` | S2 at 3 levels |
| `compressor_lz4.go` | lz4 at 3 levels |
| `compressor_xz.go` | xz/lzma2 (single level) |
| `compressor_pgzip.go` | pgzip at 3 levels |

---

## Chunk 1: Foundation

### Task 1: Initialize Go module and install dependencies

**Files:**
- Create: `go.mod`

- [ ] **Step 1: Initialize module**

```bash
cd /home/albertyw/go-compression-benchmark
go mod init github.com/albertyw/go-compression-benchmark
```

- [ ] **Step 2: Add all dependencies**

```bash
go get github.com/klauspost/compress@latest
go get github.com/andybalholm/brotli@latest
go get github.com/golang/snappy@latest
go get github.com/pierrec/lz4/v4@latest
go get github.com/ulikunitz/xz@latest
go get github.com/klauspost/pgzip@latest
```

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "Initialize Go module with compression library dependencies"
```

### Task 2: Create Compressor struct and registry

**Files:**
- Create: `compressor.go`

- [ ] **Step 1: Write compressor.go**

```go
package main

// Compressor defines a compression algorithm with its metadata and functions.
type Compressor struct {
	Name       string                       // display name, e.g. "Zstd"
	Level      string                       // compression level, e.g. "SpeedFastest"
	Library    string                       // import path, e.g. "github.com/klauspost/compress/zstd"
	PureGo     bool                         // true if no CGo dependency
	Notes      string                       // portability/edge case notes
	Compress   func([]byte) ([]byte, error) // compress input bytes
	Decompress func([]byte) ([]byte, error) // decompress back to original
}

var registry []Compressor

// Register adds a compressor to the global registry.
func Register(c Compressor) {
	registry = append(registry, c)
}

// GetAll returns all registered compressors.
func GetAll() []Compressor {
	return registry
}
```

- [ ] **Step 2: Verify it compiles**

```bash
echo 'package main; func main() {}' > main.go
go vet ./...
rm main.go
```

- [ ] **Step 3: Commit**

```bash
git add compressor.go
git commit -m "Add Compressor struct and registry"
```

### Task 3: Create benchmark runner

**Files:**
- Create: `benchmark.go`

- [ ] **Step 1: Write benchmark.go**

```go
package main

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"time"
)

const (
	warmupIterations  = 2
	measureIterations = 10
)

// BenchmarkResult holds the median metrics for one compressor on one data file.
type BenchmarkResult struct {
	Compressor       Compressor
	OriginalSize     int
	CompressedSize   int
	Ratio            float64
	CompressTime     time.Duration
	DecompressTime   time.Duration
	CompressMem      uint64
	CompressAllocs   uint64
	DecompressMem    uint64
	DecompressAllocs uint64
}

// RunBenchmarks benchmarks all registered compressors against the given data.
func RunBenchmarks(data []byte) []BenchmarkResult {
	var results []BenchmarkResult
	for _, c := range GetAll() {
		fmt.Printf("  %s (%s)...\n", c.Name, c.Level)
		result, err := benchmarkOne(c, data)
		if err != nil {
			fmt.Printf("  WARNING: skipping %s (%s): %v\n", c.Name, c.Level, err)
			continue
		}
		results = append(results, result)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].CompressedSize < results[j].CompressedSize
	})
	return results
}

func benchmarkOne(c Compressor, data []byte) (BenchmarkResult, error) {
	// Warmup
	for i := 0; i < warmupIterations; i++ {
		compressed, err := c.Compress(data)
		if err != nil {
			return BenchmarkResult{}, fmt.Errorf("compress: %w", err)
		}
		_, err = c.Decompress(compressed)
		if err != nil {
			return BenchmarkResult{}, fmt.Errorf("decompress: %w", err)
		}
	}

	// One compress to get size and verify roundtrip
	compressed, err := c.Compress(data)
	if err != nil {
		return BenchmarkResult{}, fmt.Errorf("compress: %w", err)
	}
	decompressed, err := c.Decompress(compressed)
	if err != nil {
		return BenchmarkResult{}, fmt.Errorf("decompress: %w", err)
	}
	if !bytes.Equal(data, decompressed) {
		return BenchmarkResult{}, fmt.Errorf("roundtrip mismatch")
	}

	// Timing runs
	compressTimes := make([]time.Duration, measureIterations)
	decompressTimes := make([]time.Duration, measureIterations)
	for i := 0; i < measureIterations; i++ {
		start := time.Now()
		compressed, _ = c.Compress(data)
		compressTimes[i] = time.Since(start)

		start = time.Now()
		_, _ = c.Decompress(compressed)
		decompressTimes[i] = time.Since(start)
	}

	// Memory runs
	compressMems := make([]uint64, measureIterations)
	compressAllocs := make([]uint64, measureIterations)
	decompressMems := make([]uint64, measureIterations)
	decompressAllocs := make([]uint64, measureIterations)
	for i := 0; i < measureIterations; i++ {
		var before, after runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&before)
		compressed, _ = c.Compress(data)
		runtime.ReadMemStats(&after)
		compressMems[i] = after.TotalAlloc - before.TotalAlloc
		compressAllocs[i] = after.Mallocs - before.Mallocs

		runtime.GC()
		runtime.ReadMemStats(&before)
		_, _ = c.Decompress(compressed)
		runtime.ReadMemStats(&after)
		decompressMems[i] = after.TotalAlloc - before.TotalAlloc
		decompressAllocs[i] = after.Mallocs - before.Mallocs
	}

	ratio := float64(len(data)) / float64(len(compressed))
	if len(compressed) == 0 {
		ratio = 0
	}

	return BenchmarkResult{
		Compressor:       c,
		OriginalSize:     len(data),
		CompressedSize:   len(compressed),
		Ratio:            ratio,
		CompressTime:     medianDuration(compressTimes),
		DecompressTime:   medianDuration(decompressTimes),
		CompressMem:      medianUint64(compressMems),
		CompressAllocs:   medianUint64(compressAllocs),
		DecompressMem:    medianUint64(decompressMems),
		DecompressAllocs: medianUint64(decompressAllocs),
	}, nil
}

func medianDuration(d []time.Duration) time.Duration {
	sort.Slice(d, func(i, j int) bool { return d[i] < d[j] })
	return d[len(d)/2]
}

func medianUint64(v []uint64) uint64 {
	sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
	return v[len(v)/2]
}
```

- [ ] **Step 2: Verify it compiles**

```bash
echo 'package main; func main() {}' > main.go
go vet ./...
rm main.go
```

- [ ] **Step 3: Commit**

```bash
git add benchmark.go
git commit -m "Add benchmark runner with timing and memory measurement"
```

### Task 4: Create markdown results formatter

**Files:**
- Create: `results.go`

- [ ] **Step 1: Write results.go**

```go
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
```

- [ ] **Step 2: Verify it compiles**

```bash
echo 'package main; func main() {}' > main.go
go vet ./...
rm main.go
```

- [ ] **Step 3: Commit**

```bash
git add results.go
git commit -m "Add markdown results formatter"
```

### Task 5: Create main.go entry point

**Files:**
- Create: `main.go`

- [ ] **Step 1: Write main.go**

No blank import needed — all compressor_*.go files are in `package main` and their `init()` functions run automatically.

```go
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
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add main.go
git commit -m "Add main entry point"
```

---

## Chunk 2: Compressor Implementations (Part 1)

All compressor files are in the root directory, `package main`. Each registers via `init()`.

### Task 6: Create uncompressed baseline

**Files:**
- Create: `compressor_uncompressed.go`

- [ ] **Step 1: Write compressor_uncompressed.go**

```go
package main

func init() {
	Register(Compressor{
		Name:    "(none)",
		Level:   "—",
		Library: "",
		PureGo:  false,
		Notes:   "Baseline — no compression",
		Compress: func(data []byte) ([]byte, error) {
			out := make([]byte, len(data))
			copy(out, data)
			return out, nil
		},
		Decompress: func(data []byte) ([]byte, error) {
			out := make([]byte, len(data))
			copy(out, data)
			return out, nil
		},
	})
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add compressor_uncompressed.go
git commit -m "Add uncompressed baseline compressor"
```

### Task 7: Create gzip stdlib compressor

**Files:**
- Create: `compressor_gzip_stdlib.go`

- [ ] **Step 1: Write compressor_gzip_stdlib.go**

```go
package main

import (
	"bytes"
	"compress/gzip"
	"io"
)

func init() {
	levels := []struct {
		name  string
		level int
	}{
		{"BestSpeed", gzip.BestSpeed},
		{"Default", gzip.DefaultCompression},
		{"BestCompression", gzip.BestCompression},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "Gzip stdlib",
			Level:   level.name,
			Library: "compress/gzip",
			PureGo:  true,
			Notes:   "Part of Go stdlib, guaranteed cross-platform support",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w, err := gzip.NewWriterLevel(&buf, level.level)
				if err != nil {
					return nil, err
				}
				if _, err := w.Write(data); err != nil {
					return nil, err
				}
				if err := w.Close(); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				r, err := gzip.NewReader(bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				defer r.Close()
				return io.ReadAll(r)
			},
		})
	}
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add compressor_gzip_stdlib.go
git commit -m "Add gzip stdlib compressor"
```

### Task 8: Create gzip klauspost compressor

**Files:**
- Create: `compressor_gzip_klauspost.go`

- [ ] **Step 1: Write compressor_gzip_klauspost.go**

```go
package main

import (
	"bytes"
	"io"

	kgzip "github.com/klauspost/compress/gzip"
)

func init() {
	levels := []struct {
		name  string
		level int
	}{
		{"BestSpeed", kgzip.BestSpeed},
		{"Default", kgzip.DefaultCompression},
		{"BestCompression", kgzip.BestCompression},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "Gzip klauspost",
			Level:   level.name,
			Library: "github.com/klauspost/compress/gzip",
			PureGo:  true,
			Notes:   "Drop-in stdlib replacement, assembly optimizations on amd64",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w, err := kgzip.NewWriterLevel(&buf, level.level)
				if err != nil {
					return nil, err
				}
				if _, err := w.Write(data); err != nil {
					return nil, err
				}
				if err := w.Close(); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				r, err := kgzip.NewReader(bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				defer r.Close()
				return io.ReadAll(r)
			},
		})
	}
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add compressor_gzip_klauspost.go
git commit -m "Add gzip klauspost compressor"
```

### Task 9: Create deflate klauspost compressor

**Files:**
- Create: `compressor_deflate_klauspost.go`

- [ ] **Step 1: Write compressor_deflate_klauspost.go**

```go
package main

import (
	"bytes"
	"io"

	kflate "github.com/klauspost/compress/flate"
)

func init() {
	levels := []struct {
		name  string
		level int
	}{
		{"BestSpeed", kflate.BestSpeed},
		{"Default", kflate.DefaultCompression},
		{"BestCompression", kflate.BestCompression},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "Deflate klauspost",
			Level:   level.name,
			Library: "github.com/klauspost/compress/flate",
			PureGo:  true,
			Notes:   "Drop-in stdlib replacement, assembly optimizations on amd64",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w, err := kflate.NewWriter(&buf, level.level)
				if err != nil {
					return nil, err
				}
				if _, err := w.Write(data); err != nil {
					return nil, err
				}
				if err := w.Close(); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				r := kflate.NewReader(bytes.NewReader(data))
				defer r.Close()
				return io.ReadAll(r)
			},
		})
	}
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add compressor_deflate_klauspost.go
git commit -m "Add deflate klauspost compressor"
```

### Task 10: Create zstd compressor

**Files:**
- Create: `compressor_zstd.go`

Note: Pre-allocate encoder/decoder per level to avoid measuring initialization cost in every iteration.

- [ ] **Step 1: Write compressor_zstd.go**

```go
package main

import (
	"github.com/klauspost/compress/zstd"
)

func init() {
	// Shared decoder — zstd decoding is level-independent
	dec, err := zstd.NewReader(nil)
	if err != nil {
		panic("zstd: failed to create decoder: " + err.Error())
	}

	levels := []struct {
		name  string
		level zstd.EncoderLevel
	}{
		{"SpeedFastest", zstd.SpeedFastest},
		{"SpeedDefault", zstd.SpeedDefault},
		{"SpeedBetterCompression", zstd.SpeedBetterCompression},
		{"SpeedBestCompression", zstd.SpeedBestCompression},
	}
	for _, l := range levels {
		level := l
		enc, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(level.level))
		if err != nil {
			panic("zstd: failed to create encoder: " + err.Error())
		}
		Register(Compressor{
			Name:    "Zstd",
			Level:   level.name,
			Library: "github.com/klauspost/compress/zstd",
			PureGo:  true,
			Notes:   "Pure Go, assembly on amd64, good cross-platform support",
			Compress: func(data []byte) ([]byte, error) {
				return enc.EncodeAll(data, nil), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				return dec.DecodeAll(data, nil)
			},
		})
	}
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add compressor_zstd.go
git commit -m "Add zstd compressor"
```

---

## Chunk 3: Compressor Implementations (Part 2)

### Task 11: Create brotli compressor

**Files:**
- Create: `compressor_brotli.go`

- [ ] **Step 1: Write compressor_brotli.go**

```go
package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
)

func init() {
	levels := []struct {
		name  string
		level int
	}{
		{"Fast", 1},
		{"Default", 6},
		{"Best", 11},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "Brotli",
			Level:   fmt.Sprintf("%s (%d)", level.name, level.level),
			Library: "github.com/andybalholm/brotli",
			PureGo:  true,
			Notes:   "Pure Go, fully portable",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w := brotli.NewWriterLevel(&buf, level.level)
				if _, err := w.Write(data); err != nil {
					return nil, err
				}
				if err := w.Close(); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				r := brotli.NewReader(bytes.NewReader(data))
				return io.ReadAll(r)
			},
		})
	}
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add compressor_brotli.go
git commit -m "Add brotli compressor"
```

### Task 12: Create snappy compressor

**Files:**
- Create: `compressor_snappy.go`

- [ ] **Step 1: Write compressor_snappy.go**

```go
package main

import (
	"github.com/golang/snappy"
)

func init() {
	Register(Compressor{
		Name:    "Snappy",
		Level:   "—",
		Library: "github.com/golang/snappy",
		PureGo:  true,
		Notes:   "Pure Go, assembly on amd64",
		Compress: func(data []byte) ([]byte, error) {
			return snappy.Encode(nil, data), nil
		},
		Decompress: func(data []byte) ([]byte, error) {
			return snappy.Decode(nil, data)
		},
	})
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add compressor_snappy.go
git commit -m "Add snappy compressor"
```

### Task 13: Create S2 compressor

**Files:**
- Create: `compressor_s2.go`

- [ ] **Step 1: Write compressor_s2.go**

```go
package main

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/s2"
)

func init() {
	type s2Level struct {
		name    string
		options []s2.WriterOption
	}
	levels := []s2Level{
		{"Default", nil},
		{"Better", []s2.WriterOption{s2.WriterBetterCompression()}},
		{"Best", []s2.WriterOption{s2.WriterBestCompression()}},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "S2",
			Level:   level.name,
			Library: "github.com/klauspost/compress/s2",
			PureGo:  true,
			Notes:   "Snappy-compatible decoding; S2-compressed data cannot be decoded by Snappy decoders",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w := s2.NewWriter(&buf, level.options...)
				if _, err := w.Write(data); err != nil {
					return nil, err
				}
				if err := w.Close(); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				r := s2.NewReader(bytes.NewReader(data))
				return io.ReadAll(r)
			},
		})
	}
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add compressor_s2.go
git commit -m "Add S2 compressor"
```

### Task 14: Create LZ4 compressor

**Files:**
- Create: `compressor_lz4.go`

- [ ] **Step 1: Write compressor_lz4.go**

```go
package main

import (
	"bytes"
	"io"

	"github.com/pierrec/lz4/v4"
)

func init() {
	levels := []struct {
		name  string
		level lz4.CompressionLevel
	}{
		{"Fast", lz4.Fast},
		{"Level5", lz4.Level5},
		{"Level9", lz4.Level9},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "LZ4",
			Level:   level.name,
			Library: "github.com/pierrec/lz4/v4",
			PureGo:  true,
			Notes:   "Pure Go, fully portable",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w := lz4.NewWriter(&buf)
				if err := w.Apply(lz4.CompressionLevelOption(level.level)); err != nil {
					return nil, err
				}
				if _, err := w.Write(data); err != nil {
					return nil, err
				}
				if err := w.Close(); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				r := lz4.NewReader(bytes.NewReader(data))
				return io.ReadAll(r)
			},
		})
	}
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add compressor_lz4.go
git commit -m "Add LZ4 compressor"
```

### Task 15: Create XZ compressor

**Files:**
- Create: `compressor_xz.go`

- [ ] **Step 1: Write compressor_xz.go**

```go
package main

import (
	"bytes"
	"io"

	"github.com/ulikunitz/xz"
)

func init() {
	Register(Compressor{
		Name:    "XZ/LZMA2",
		Level:   "—",
		Library: "github.com/ulikunitz/xz",
		PureGo:  true,
		Notes:   "Pure Go, fully portable but significantly slower than C implementation",
		Compress: func(data []byte) ([]byte, error) {
			var buf bytes.Buffer
			w, err := xz.NewWriter(&buf)
			if err != nil {
				return nil, err
			}
			if _, err := w.Write(data); err != nil {
				return nil, err
			}
			if err := w.Close(); err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		},
		Decompress: func(data []byte) ([]byte, error) {
			r, err := xz.NewReader(bytes.NewReader(data))
			if err != nil {
				return nil, err
			}
			return io.ReadAll(r)
		},
	})
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add compressor_xz.go
git commit -m "Add XZ/LZMA2 compressor"
```

### Task 16: Create pgzip compressor

**Files:**
- Create: `compressor_pgzip.go`

- [ ] **Step 1: Write compressor_pgzip.go**

```go
package main

import (
	"bytes"
	"io"

	"github.com/klauspost/pgzip"
)

func init() {
	levels := []struct {
		name  string
		level int
	}{
		{"BestSpeed", pgzip.BestSpeed},
		{"Default", pgzip.DefaultCompression},
		{"BestCompression", pgzip.BestCompression},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "pgzip",
			Level:   level.name,
			Library: "github.com/klauspost/pgzip",
			PureGo:  true,
			Notes:   "Parallel gzip — performance scales with CPU cores, results vary by machine",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w, err := pgzip.NewWriterLevel(&buf, level.level)
				if err != nil {
					return nil, err
				}
				if _, err := w.Write(data); err != nil {
					return nil, err
				}
				if err := w.Close(); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				r, err := pgzip.NewReader(bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				defer r.Close()
				return io.ReadAll(r)
			},
		})
	}
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go vet ./...
```

- [ ] **Step 3: Commit**

```bash
git add compressor_pgzip.go
git commit -m "Add pgzip compressor"
```

---

## Chunk 4: Integration and Final Run

### Task 17: Build and run the full benchmark

- [ ] **Step 1: Run go mod tidy**

```bash
go mod tidy
```

- [ ] **Step 2: Build**

```bash
go build -o benchmark .
```

- [ ] **Step 3: Run**

```bash
./benchmark
```

Expected: prints progress for each data file and compressor, writes RESULTS.md.

- [ ] **Step 4: Verify RESULTS.md was generated and looks correct**

Check that the file has two tables (data.h3 and data_mock.h3), all compressors appear, ratios are reasonable, and no errors.

- [ ] **Step 5: Commit results**

```bash
git add RESULTS.md
git commit -m "Add initial benchmark results"
```

### Task 18: Update CLAUDE.md and final cleanup

- [ ] **Step 1: Update CLAUDE.md**

Update the Architecture section to reflect the flat `package main` structure (no `compressors/` subdirectory). Compressor files are named `compressor_*.go` in the root directory.

- [ ] **Step 2: Run go vet**

```bash
go vet ./...
```

- [ ] **Step 3: Check for any issues in RESULTS.md**

If data_mock.h3 results are uninformative (sub-microsecond timings, ratios near 1.0), consider noting this in the output.

- [ ] **Step 4: Final commit**

```bash
git add -A
git commit -m "Update CLAUDE.md and final cleanup"
```
