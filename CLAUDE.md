# Go Compression Benchmark

## Project Overview
Benchmarks Go compression libraries for use in [localtimezone](https://github.com/albertyw/localtimezone). Compresses/decompresses H3 timezone binary data and outputs results as markdown tables in `RESULTS.md`.

## Build & Run
```bash
go mod tidy          # install dependencies
go build -o benchmark .  # build
./benchmark          # run benchmarks, writes RESULTS.md
go vet ./...         # lint
```

## Architecture
- `main.go` — entry point: loads data files, runs benchmarks, writes RESULTS.md
- `compressor.go` — Compressor struct (Name, Level, Library, etc.) + registry (Register/GetAll)
- `compressor_*.go` — one file per library in root package main, registers via init()
- `benchmark.go` — benchmark runner with warmup, timing, and memory measurement passes
- `results.go` — markdown table formatter

## Conventions
- Module: `github.com/albertyw/go-compression-benchmark`
- All compression libraries must be pure Go (no CGo)
- Each compressor file registers variants via `init()` calling `Register()`
- Benchmark runner separates timing runs from memory measurement runs
- Test data: `data.h3` (real, ~8.5MB), `data_mock.h3` (mock, ~1.2KB)
