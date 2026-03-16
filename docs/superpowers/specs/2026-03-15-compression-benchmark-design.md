# Go Compression Benchmark Design

## Purpose

Benchmark well-supported Go compression libraries to determine the best option for [localtimezone](https://github.com/albertyw/localtimezone), which needs to compress/decompress H3 timezone binary data. The output is a `RESULTS.md` file containing markdown tables comparing compression efficiency, speed, memory usage, and allocation counts.

## Libraries

All libraries are pure Go (no CGo), ensuring full cross-platform portability.

| Library | Package | Levels |
|---------|---------|--------|
| Uncompressed | (baseline) | Single |
| Gzip (stdlib) | `compress/gzip` | BestSpeed (1), Default (6), BestCompression (9) |
| Gzip (klauspost) | `github.com/klauspost/compress/gzip` | BestSpeed (1), Default (6), BestCompression (9) |
| Deflate (klauspost) | `github.com/klauspost/compress/flate` | BestSpeed (1), Default (6), BestCompression (9) |
| Zstd | `github.com/klauspost/compress/zstd` | SpeedFastest, SpeedDefault, SpeedBetterCompression, SpeedBestCompression |
| Brotli | `github.com/andybalholm/brotli` | 1 (fast), 6 (default), 11 (best) |
| Snappy | `github.com/golang/snappy` | Single |
| S2 | `github.com/klauspost/compress/s2` | Default (no option), Better (WriterBetterCompression), Best (WriterBestCompression) |
| LZ4 | `github.com/pierrec/lz4/v4` | Fast, Level 5 (mid), Level 9 (high) |
| XZ/LZMA2 | `github.com/ulikunitz/xz` | Single |
| pgzip | `github.com/klauspost/pgzip` | BestSpeed (1), Default (6), BestCompression (9) |

## Architecture

```
main.go              — entry point: loads data, runs benchmarks, writes RESULTS.md
compressor.go        — Compressor struct definition + registry (Register/GetAll)
compressors/         — one file per library
  uncompressed.go
  gzip_stdlib.go
  gzip_klauspost.go
  deflate_klauspost.go
  zstd.go
  brotli.go
  snappy.go
  s2.go
  lz4.go
  xz.go
  pgzip.go
results.go           — markdown table formatter, writes RESULTS.md
go.mod / go.sum      — module: github.com/albertyw/go-compression-benchmark
data.h3              — real test data (~8.5 MB)
data_mock.h3         — small test data (~1.2 KB)
RESULTS.md           — generated output
```

### Compressor struct

```go
type Compressor struct {
    Name       string                         // e.g. "Zstd (fastest)"
    Library    string                         // e.g. "github.com/klauspost/compress/zstd"
    PureGo     bool                           // all selected libraries are pure Go; kept for future extensibility
    Notes      string                         // portability/edge case notes
    Compress   func([]byte) ([]byte, error)   // compress input bytes
    Decompress func([]byte) ([]byte, error)   // decompress back to original
}
```

Using a struct with function fields rather than an interface. Each compressor file creates closures that capture compression level configuration and calls `Register()` via `init()`.

### Registry

```go
var registry []Compressor

func Register(c Compressor) {
    registry = append(registry, c)
}

func GetAll() []Compressor {
    return registry
}
```

## Benchmark Runner

For each compressor and each data file:

1. Run 2 warmup iterations (discarded) to warm CPU caches and trigger lazy initialization
2. Run 10 measured iterations, separated into timing runs and memory runs:
   - **Timing runs:** Use `time.Now()` (monotonic clock via `time.Since`) with no `ReadMemStats` calls to avoid stop-the-world pauses distorting timing
   - **Memory runs:** Call `runtime.GC()`, then `runtime.ReadMemStats` before and after compress/decompress to measure `TotalAlloc` and `Mallocs` deltas
3. Verify decompressed output matches original input
4. Take median of 10 iterations for all metrics
5. If a compressor returns an error, log a warning and skip that entry (do not abort)

### Metrics

| Metric | Measurement |
|--------|-------------|
| Compressed size (bytes) | `len(compressed)` |
| Compression ratio | `float64(len(original)) / float64(len(compressed))`, displayed as "2.34x" |
| Compression time | `time.Since(start)` (monotonic clock, timing-only runs) |
| Decompression time | Same |
| Compression memory (bytes) | `MemStats.TotalAlloc` delta (memory-only runs after `runtime.GC()`) |
| Compression allocations | `MemStats.Mallocs` delta (memory-only runs) |
| Decompression memory (bytes) | `MemStats.TotalAlloc` delta (memory-only runs) |
| Decompression allocations | `MemStats.Mallocs` delta (memory-only runs) |

## Output Format

`RESULTS.md` contains one table per data file:

```markdown
# Compression Benchmark Results

## data.h3 (8,966,666 bytes)

| Library | Level | Compressed Size | Ratio | Compress Time | Decompress Time | Compress Mem | Compress Allocs | Decompress Mem | Decompress Allocs | Pure Go | Notes |
|---------|-------|----------------|-------|--------------|----------------|-------------|----------------|---------------|-------------------|---------|-------|
| (none) | — | 8,966,666 B | 1.00x | — | — | 0 B | 0 | 0 B | 0 | — | Baseline |
| Zstd | fastest | ... | ... | ... | ... | ... | ... | ... | ... | Yes | ... |
| ... | ... | ... | ... | ... | ... | ... | ... | ... | ... | ... | ... |

## data_mock.h3 (1,253 bytes)

| Library | Level | ... |
| ... |

## Notes

- All libraries are pure Go with no CGo dependencies
- Benchmarks run 2 warmup + 10 measured iterations per compressor; median values reported
- Times measured via monotonic clock (time.Since) in timing-only runs
- Memory/allocs measured via runtime.MemStats deltas in separate memory-only runs (after runtime.GC())
- Rows sorted by compressed size ascending within each table
- [Additional manually curated portability notes per library]
```

## Data Files

- `data.h3` (~8.5 MB) — real H3 timezone binary data from localtimezone, primary benchmark target
- `data_mock.h3` (~1.2 KB) — small mock data for quick validation; results may be discarded if uninformative

## Edge Cases & Portability Notes (manually curated per library)

These are embedded in each compressor's `Notes` field:

- **Gzip (stdlib):** Part of Go stdlib, guaranteed cross-platform support
- **Gzip/Deflate (klauspost):** Drop-in stdlib replacement, uses assembly optimizations on amd64
- **Zstd:** Pure Go, assembly on amd64, good cross-platform support
- **Brotli:** Pure Go, fully portable
- **Snappy:** Pure Go, assembly on amd64
- **S2:** Snappy-compatible decoding but S2-compressed data cannot be decoded by Snappy decoders
- **LZ4:** Pure Go, fully portable
- **XZ/LZMA2:** Pure Go, fully portable but significantly slower than C implementation
- **pgzip:** Parallel — performance scales with CPU cores, results vary by machine
