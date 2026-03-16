# Go Compression Benchmark

Comparing performance of Go compression libraries against real-world binary data ([H3 timezone cells](https://github.com/albertyw/localtimezone)). Measures compression ratio, speed, memory usage, and allocations across multiple compression levels.

## Libraries Tested

| Library | Package | Pure Go |
|---------|---------|---------|
| Gzip (stdlib) | `compress/gzip` | Yes |
| Gzip (klauspost) | `github.com/klauspost/compress/gzip` | Yes |
| Deflate (klauspost) | `github.com/klauspost/compress/flate` | Yes |
| Zstd | `github.com/klauspost/compress/zstd` | Yes |
| Brotli | `github.com/andybalholm/brotli` | Yes |
| Snappy | `github.com/golang/snappy` | Yes |
| S2 | `github.com/klauspost/compress/s2` | Yes |
| LZ4 | `github.com/pierrec/lz4/v4` | Yes |
| XZ/LZMA2 | `github.com/ulikunitz/xz` | Yes |
| pgzip | `github.com/klauspost/pgzip` | Yes |

All libraries are pure Go with no CGo dependencies, ensuring full cross-platform portability.

## Results

See [RESULTS.md](RESULTS.md) for full benchmark tables. Highlights on the 8.5 MB test file:

| Library | Level | Ratio | Compress | Decompress |
|---------|-------|-------|----------|------------|
| XZ/LZMA2 | default | 10.33x | 449ms | 97ms |
| Zstd | SpeedFastest | 5.78x | 34ms | 13ms |
| Zstd | SpeedBestCompression | 6.67x | 317ms | 14ms |
| Brotli | Best (11) | 7.79x | 10.1s | 25ms |
| pgzip | BestSpeed | 5.16x | 5.6ms | 14ms |
| Snappy | default | 2.19x | 8ms | 6ms |
| S2 | Better | 2.40x | 5.9ms | 6ms |

## Build & Run

```bash
go build -o benchmark .
./benchmark
```

This reads the test data files and writes results to `RESULTS.md`.

## Test Data

- `data.h3` (~8.5 MB) -- real H3 timezone binary data from [localtimezone](https://github.com/albertyw/localtimezone)
- `data_mock.h3` (~1.2 KB) -- small mock data for quick validation

## Architecture

- `main.go` -- entry point: loads data files, runs benchmarks, writes RESULTS.md
- `compressor.go` -- `Compressor` struct (Name, Level, Library, compress/decompress functions) and registry
- `compressor_*.go` -- one file per library, registers variants via `init()`
- `benchmark.go` -- benchmark runner with warmup, separate timing and memory measurement passes, median calculation
- `results.go` -- markdown table formatter

## Methodology

- 2 warmup iterations (discarded) to warm CPU caches and trigger lazy initialization
- 10 measured iterations per compressor, with median reported
- Timing and memory measured in separate passes to avoid `ReadMemStats` stop-the-world pauses distorting timing
- Roundtrip verification ensures decompressed output matches original input
- "Ratio" = original size / compressed size (higher is better)
- Memory figures are approximate; `runtime.MemStats` uses process-global counters, so allocations from concurrent goroutines (e.g. pgzip workers, GC) may be included

## Adding a Compressor

Create a new `compressor_name.go` file in the root directory:

```go
package main

func init() {
    Register(Compressor{
        Name:    "MyLib",
        Level:   "Default",
        Library: "github.com/example/mylib",
        PureGo:  true,
        Notes:   "Description of portability characteristics",
        Compress: func(data []byte) ([]byte, error) {
            // compress data and return
        },
        Decompress: func(data []byte) ([]byte, error) {
            // decompress data and return
        },
    })
}
```

## License

MIT
