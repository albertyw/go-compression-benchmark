# Compression Benchmark Results

## Environment

- **CPU:** AMD Ryzen 9 7900X 12-Core Processor
- **Architecture:** amd64
- **OS:** linux
- **Go:** go1.26.1

## data_mock.h3 (1,253 bytes)

### Compression

| Library | Level | Compressed Size | Ratio | Time | Memory | Allocs | Pure Go | Notes |
|---------|-------|----------------|-------|------|--------|--------|---------|-------|
| [Brotli](https://github.com/andybalholm/brotli) | Best (11) | 142 B | 8.82x | 1.27ms | 33.01 MB | 61 | Yes | Pure Go, fully portable |
| [Brotli](https://github.com/andybalholm/brotli) | Default (6) | 156 B | 8.03x | 106.3µs | 2.09 MB | 40 | Yes | Pure Go, fully portable |
| [Zstd](https://github.com/klauspost/compress) | SpeedFastest | 191 B | 6.56x | 14.3µs | 1.25 KB | 1 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| [Zstd](https://github.com/klauspost/compress) | SpeedDefault | 191 B | 6.56x | 16.6µs | 1.25 KB | 1 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| [Zstd](https://github.com/klauspost/compress) | SpeedBetterCompression | 191 B | 6.56x | 27.9µs | 1.25 KB | 1 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| [Zstd](https://github.com/klauspost/compress) | SpeedBestCompression | 194 B | 6.46x | 101.6µs | 1.25 KB | 1 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| [Deflate klauspost](https://github.com/klauspost/compress) | BestCompression | 232 B | 5.40x | 65.2µs | 1.09 MB | 16 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [Deflate klauspost](https://github.com/klauspost/compress) | Default | 242 B | 5.18x | 42.4µs | 1.03 MB | 15 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [Brotli](https://github.com/andybalholm/brotli) | Fast (1) | 242 B | 5.18x | 9.7µs | 44.16 KB | 12 | Yes | Pure Go, fully portable |
| [Deflate klauspost](https://github.com/klauspost/compress) | BestSpeed | 248 B | 5.05x | 48.0µs | 794.92 KB | 15 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [Gzip klauspost](https://github.com/klauspost/compress) | BestCompression | 250 B | 5.01x | 69.3µs | 1.09 MB | 18 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | Default | 252 B | 4.97x | 50.5µs | 795.59 KB | 21 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | BestCompression | 252 B | 4.97x | 47.8µs | 795.59 KB | 21 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| [pgzip](https://github.com/klauspost/pgzip) | BestCompression | 255 B | 4.91x | 233.9µs | 4.29 MB | 46 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| [Gzip klauspost](https://github.com/klauspost/compress) | Default | 260 B | 4.82x | 41.7µs | 1.03 MB | 18 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | BestSpeed | 263 B | 4.76x | 47.9µs | 1.15 MB | 23 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| [pgzip](https://github.com/klauspost/pgzip) | Default | 265 B | 4.73x | 207.3µs | 4.22 MB | 45 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| [Gzip klauspost](https://github.com/klauspost/compress) | BestSpeed | 266 B | 4.71x | 75.7µs | 795.64 KB | 18 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [pgzip](https://github.com/klauspost/pgzip) | BestSpeed | 271 B | 4.62x | 195.6µs | 3.97 MB | 45 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| [XZ/LZMA2](https://github.com/ulikunitz/xz) | — | 272 B | 4.61x | 2.10ms | 48.13 MB | 414 | Yes | Pure Go, fully portable but significantly slower than C implementation |
| [Snappy](https://github.com/golang/snappy) | — | 525 B | 2.39x | 0.9µs | 1.50 KB | 1 | Yes | Pure Go, assembly on amd64 |
| [LZ4](https://github.com/pierrec/lz4) | Level9 | 556 B | 2.25x | 11.3µs | 5.02 MB | 14 | Yes | Pure Go, fully portable |
| [LZ4](https://github.com/pierrec/lz4) | Level5 | 556 B | 2.25x | 11.2µs | 5.02 MB | 14 | Yes | Pure Go, fully portable |
| [S2](https://github.com/klauspost/compress) | Best | 564 B | 2.22x | 310.4µs | 6.51 MB | 18 | Yes | S2 can decode Snappy data, but Snappy decoders cannot decode S2 output |
| [S2](https://github.com/klauspost/compress) | Better | 565 B | 2.22x | 90.6µs | 2.03 MB | 19 | Yes | S2 can decode Snappy data, but Snappy decoders cannot decode S2 output |
| [LZ4](https://github.com/pierrec/lz4) | Fast | 568 B | 2.21x | 1.8µs | 4.14 MB | 14 | Yes | Pure Go, fully portable |
| [S2](https://github.com/klauspost/compress) | Default | 581 B | 2.16x | 86.3µs | 2.02 MB | 19 | Yes | S2 can decode Snappy data, but Snappy decoders cannot decode S2 output |
| (none) | — | 1.22 KB | 1.00x | 0.1µs | 1.25 KB | 1 | — | Baseline — no compression |

### Decompression

| Library | Level | Time | Memory | Allocs |
|---------|-------|------|--------|--------|
| [Brotli](https://github.com/andybalholm/brotli) | Best (11) | 12.2µs | 73.70 KB | 19 |
| [Brotli](https://github.com/andybalholm/brotli) | Default (6) | 11.2µs | 65.76 KB | 19 |
| [Zstd](https://github.com/klauspost/compress) | SpeedFastest | 5.0µs | 7.70 KB | 7 |
| [Zstd](https://github.com/klauspost/compress) | SpeedDefault | 5.1µs | 7.70 KB | 7 |
| [Zstd](https://github.com/klauspost/compress) | SpeedBetterCompression | 4.4µs | 7.70 KB | 7 |
| [Zstd](https://github.com/klauspost/compress) | SpeedBestCompression | 5.8µs | 7.70 KB | 7 |
| [Deflate klauspost](https://github.com/klauspost/compress) | BestCompression | 5.7µs | 40.08 KB | 12 |
| [Deflate klauspost](https://github.com/klauspost/compress) | Default | 6.4µs | 40.08 KB | 12 |
| [Brotli](https://github.com/andybalholm/brotli) | Fast (1) | 9.6µs | 65.76 KB | 19 |
| [Deflate klauspost](https://github.com/klauspost/compress) | BestSpeed | 6.4µs | 40.08 KB | 12 |
| [Gzip klauspost](https://github.com/klauspost/compress) | BestCompression | 6.4µs | 40.77 KB | 13 |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | Default | 6.5µs | 43.20 KB | 11 |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | BestCompression | 6.2µs | 43.20 KB | 11 |
| [pgzip](https://github.com/klauspost/pgzip) | BestCompression | 196.3µs | 4.04 MB | 29 |
| [Gzip klauspost](https://github.com/klauspost/compress) | Default | 5.5µs | 40.77 KB | 13 |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | BestSpeed | 8.2µs | 43.20 KB | 11 |
| [pgzip](https://github.com/klauspost/pgzip) | Default | 177.7µs | 4.04 MB | 29 |
| [Gzip klauspost](https://github.com/klauspost/compress) | BestSpeed | 8.5µs | 40.77 KB | 13 |
| [pgzip](https://github.com/klauspost/pgzip) | BestSpeed | 189.5µs | 4.04 MB | 29 |
| [XZ/LZMA2](https://github.com/ulikunitz/xz) | — | 45.6µs | 8.04 MB | 261 |
| [Snappy](https://github.com/golang/snappy) | — | 0.6µs | 1.25 KB | 1 |
| [LZ4](https://github.com/pierrec/lz4) | Level9 | 1.4µs | 4.01 MB | 14 |
| [LZ4](https://github.com/pierrec/lz4) | Level5 | 1.5µs | 4.01 MB | 14 |
| [S2](https://github.com/klauspost/compress) | Best | 35.1µs | 1.01 MB | 9 |
| [S2](https://github.com/klauspost/compress) | Better | 39.6µs | 1.01 MB | 9 |
| [LZ4](https://github.com/pierrec/lz4) | Fast | 1.6µs | 4.01 MB | 14 |
| [S2](https://github.com/klauspost/compress) | Default | 36.8µs | 1.01 MB | 9 |
| (none) | — | 0.1µs | 1.25 KB | 1 |

## Notes

- All libraries are pure Go with no CGo dependencies
- Benchmarks run 2 warmup + 10 measured iterations per compressor; median values reported
- Times measured via monotonic clock (time.Since) in timing-only runs
- Memory/allocs measured via runtime.MemStats deltas in separate memory-only runs (after runtime.GC())
- Rows sorted by compressed size ascending within each table
- Ratio = original size / compressed size (higher is better)
- Memory figures are approximate; runtime.MemStats uses process-global counters, so allocations from concurrent goroutines may be included
