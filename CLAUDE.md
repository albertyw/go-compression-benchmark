# Go Compression Benchmark

## Project Overview
Benchmarks Go compression libraries for use in [localtimezone](https://github.com/albertyw/localtimezone). Compresses/decompresses H3 timezone binary data and outputs results as markdown tables in `results/`.

## Build & Run
```bash
make test       # run tests
make vet        # lint
make benchmark  # build and run benchmarks
go mod tidy     # install dependencies
```

## Architecture
- `main.go` — entry point: discovers files in `sample_data/`, runs benchmarks, writes `results/`
- `compressors/compressor.go` — Compressor struct (Name, Level, Library, etc.) + registry (Register/GetAll/Registry)
- `compressors/compressor_*.go` — one file per library in `package compressors`, registers via `init()`
- `benchmark.go` — benchmark runner with warmup, timing, and memory measurement passes
- `results.go` — markdown table formatter
- `RESULTS.md` — summary blog post / analysis of benchmark results

## Conventions
- Module: `github.com/albertyw/go-compression-benchmark`
- All compression libraries must be pure Go (no CGo)
- Each compressor file registers variants via `init()` calling `Register()`
- Benchmark runner separates timing runs from memory measurement runs
- Test data lives in `sample_data/`; benchmarker auto-discovers all files there
- Results are written to `results/<datafile>-<date>-<arch>-<os>.md`
