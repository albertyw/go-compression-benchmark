# Go Compression Benchmark

[![Build Status](https://drone.albertyw.com/api/badges/albertyw/go-compression-benchmark/status.svg)](https://drone.albertyw.com/albertyw/go-compression-benchmark)

Read a summary at https://www.albertyw.com/note/go-compression-benchmarks

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

See [RESULTS.md](RESULTS.md) for the summary and analysis. Raw per-run tables are in [`results/`](results/).

## Build & Run

```bash
make test       # run tests
make benchmark  # build and run benchmarks, writes to results/
```

Or manually:

```bash
go mod tidy
go build -o benchmark .
./benchmark
```

Benchmark output is written to `results/<datafile>-<date>-<arch>-<os>.md`.

## Test Data

Sample data files live in `sample_data/`. The benchmarker auto-discovers all files in that directory.

- `data.h3` (~8.5 MB) — real H3 timezone binary data from [localtimezone](https://github.com/albertyw/localtimezone)
- `data_mock.h3` (~1.2 KB) — small mock data for quick validation

## Architecture

- `main.go` — entry point: discovers files in `sample_data/`, runs benchmarks, writes `results/`
- `compressors/compressor.go` — `Compressor` struct and registry (`Register`/`GetAll`)
- `compressors/compressor_*.go` — one file per library, registers variants via `init()`
- `benchmark.go` — benchmark runner with warmup, separate timing and memory passes, median calculation
- `results.go` — markdown table formatter

## Methodology

- 2 warmup iterations (discarded) to warm CPU caches and trigger lazy initialization
- 10 measured iterations per compressor, with median reported
- Timing and memory measured in separate passes to avoid `ReadMemStats` stop-the-world pauses distorting timing
- Roundtrip verification ensures decompressed output matches original input
- "Ratio" = original size / compressed size (higher is better)
- Memory figures are approximate; `runtime.MemStats` uses process-global counters, so allocations from concurrent goroutines (e.g. pgzip workers, GC) may be included

## Adding a Compressor

Create a new `compressors/compressor_name.go` file:

```go
package compressors

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
