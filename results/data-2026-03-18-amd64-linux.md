# Compression Benchmark Results

## Environment

- **CPU:** AMD Ryzen 9 7900X 12-Core Processor
- **Architecture:** amd64
- **OS:** linux
- **Go:** go1.26.1

## data.h3 (8,966,666 bytes)

### Compression

| Library | Level | Compressed Size | Ratio | Time | Memory | Allocs | Pure Go | Notes |
|---------|-------|----------------|-------|------|--------|--------|---------|-------|
| [XZ/LZMA2](https://github.com/ulikunitz/xz) | — | 847.99 KB | 10.33x | 450.00ms | 64.76 MB | 932,862 | Yes | Pure Go, fully portable but significantly slower than C implementation |
| [Brotli](https://github.com/andybalholm/brotli) | Best (11) | 1.10 MB | 7.79x | 10.65s | 802.51 MB | 372 | Yes | Pure Go, fully portable |
| [Brotli](https://github.com/andybalholm/brotli) | Default (6) | 1.20 MB | 7.15x | 302.47ms | 338.97 MB | 134 | Yes | Pure Go, fully portable |
| [Zstd](https://github.com/klauspost/compress) | SpeedBestCompression | 1.28 MB | 6.67x | 356.94ms | 23.22 MB | 63 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| [Zstd](https://github.com/klauspost/compress) | SpeedBetterCompression | 1.47 MB | 5.83x | 70.43ms | 26.16 MB | 66 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| [Zstd](https://github.com/klauspost/compress) | SpeedFastest | 1.48 MB | 5.78x | 31.39ms | 16.55 MB | 63 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| [Zstd](https://github.com/klauspost/compress) | SpeedDefault | 1.52 MB | 5.64x | 37.49ms | 26.16 MB | 66 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| [Brotli](https://github.com/andybalholm/brotli) | Fast (1) | 1.60 MB | 5.36x | 41.22ms | 14.90 MB | 13 | Yes | Pure Go, fully portable |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | Default | 1.60 MB | 5.36x | 196.95ms | 4.78 MB | 33 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | BestCompression | 1.60 MB | 5.35x | 291.88ms | 4.78 MB | 33 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| [pgzip](https://github.com/klauspost/pgzip) | BestCompression | 1.60 MB | 5.35x | 47.85ms | 45.64 MB | 301 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| [Deflate klauspost](https://github.com/klauspost/compress) | BestCompression | 1.60 MB | 5.35x | 293.16ms | 5.09 MB | 29 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [Gzip klauspost](https://github.com/klauspost/compress) | BestCompression | 1.60 MB | 5.35x | 298.61ms | 5.09 MB | 31 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [Deflate klauspost](https://github.com/klauspost/compress) | Default | 1.62 MB | 5.26x | 34.24ms | 5.03 MB | 28 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [Gzip klauspost](https://github.com/klauspost/compress) | Default | 1.62 MB | 5.26x | 35.15ms | 5.03 MB | 30 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [pgzip](https://github.com/klauspost/pgzip) | Default | 1.62 MB | 5.26x | 8.03ms | 42.06 MB | 290 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| [Deflate klauspost](https://github.com/klauspost/compress) | BestSpeed | 1.66 MB | 5.16x | 18.95ms | 4.78 MB | 28 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [Gzip klauspost](https://github.com/klauspost/compress) | BestSpeed | 1.66 MB | 5.16x | 19.36ms | 4.78 MB | 30 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| [pgzip](https://github.com/klauspost/pgzip) | BestSpeed | 1.66 MB | 5.16x | 5.62ms | 40.04 MB | 295 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | BestSpeed | 1.91 MB | 4.47x | 34.12ms | 5.15 MB | 35 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| [LZ4](https://github.com/pierrec/lz4) | Level9 | 3.52 MB | 2.43x | 107.20ms | 20.61 MB | 21 | Yes | Pure Go, fully portable |
| [LZ4](https://github.com/pierrec/lz4) | Level5 | 3.52 MB | 2.43x | 107.06ms | 20.61 MB | 21 | Yes | Pure Go, fully portable |
| [LZ4](https://github.com/pierrec/lz4) | Fast | 3.55 MB | 2.41x | 15.95ms | 20.09 MB | 23 | Yes | Pure Go, fully portable |
| [S2](https://github.com/klauspost/compress) | Better | 3.56 MB | 2.40x | 4.88ms | 27.31 MB | 138 | Yes | S2 can decode Snappy data, but Snappy decoders cannot decode S2 output |
| [S2](https://github.com/klauspost/compress) | Best | 3.57 MB | 2.40x | 52.96ms | 72.85 MB | 156 | Yes | S2 can decode Snappy data, but Snappy decoders cannot decode S2 output |
| [Snappy](https://github.com/golang/snappy) | — | 3.91 MB | 2.19x | 7.63ms | 9.98 MB | 1 | Yes | Pure Go, assembly on amd64 |
| [S2](https://github.com/klauspost/compress) | Default | 4.14 MB | 2.07x | 5.25ms | 32.63 MB | 148 | Yes | S2 can decode Snappy data, but Snappy decoders cannot decode S2 output |
| (none) | — | 8.55 MB | 1.00x | 480.5µs | 8.55 MB | 1 | — | Baseline — no compression |

### Decompression

| Library | Level | Time | Memory | Allocs |
|---------|-------|------|--------|--------|
| [XZ/LZMA2](https://github.com/ulikunitz/xz) | — | 96.65ms | 43.18 MB | 931,690 |
| [Brotli](https://github.com/andybalholm/brotli) | Best (11) | 35.98ms | 26.30 MB | 54 |
| [Brotli](https://github.com/andybalholm/brotli) | Default (6) | 27.35ms | 25.08 MB | 54 |
| [Zstd](https://github.com/klauspost/compress) | SpeedBestCompression | 13.23ms | 8.56 MB | 7 |
| [Zstd](https://github.com/klauspost/compress) | SpeedBetterCompression | 14.32ms | 8.56 MB | 7 |
| [Zstd](https://github.com/klauspost/compress) | SpeedFastest | 12.80ms | 8.56 MB | 7 |
| [Zstd](https://github.com/klauspost/compress) | SpeedDefault | 14.35ms | 8.56 MB | 7 |
| [Brotli](https://github.com/andybalholm/brotli) | Fast (1) | 39.29ms | 29.53 MB | 661 |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | Default | 24.39ms | 21.68 MB | 3,623 |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | BestCompression | 21.78ms | 21.68 MB | 3,558 |
| [pgzip](https://github.com/klauspost/pgzip) | BestCompression | 13.80ms | 24.99 MB | 196 |
| [Deflate klauspost](https://github.com/klauspost/compress) | BestCompression | 13.96ms | 20.99 MB | 171 |
| [Gzip klauspost](https://github.com/klauspost/compress) | BestCompression | 14.80ms | 20.99 MB | 171 |
| [Deflate klauspost](https://github.com/klauspost/compress) | Default | 14.49ms | 21.02 MB | 472 |
| [Gzip klauspost](https://github.com/klauspost/compress) | Default | 15.08ms | 21.02 MB | 473 |
| [pgzip](https://github.com/klauspost/pgzip) | Default | 15.78ms | 25.02 MB | 491 |
| [Deflate klauspost](https://github.com/klauspost/compress) | BestSpeed | 16.41ms | 21.03 MB | 536 |
| [Gzip klauspost](https://github.com/klauspost/compress) | BestSpeed | 13.71ms | 21.03 MB | 538 |
| [pgzip](https://github.com/klauspost/pgzip) | BestSpeed | 16.80ms | 25.03 MB | 564 |
| [Gzip stdlib](https://pkg.go.dev/compress/gzip) | BestSpeed | 24.84ms | 21.65 MB | 4,284 |
| [LZ4](https://github.com/pierrec/lz4) | Level9 | 6.80ms | 24.94 MB | 40 |
| [LZ4](https://github.com/pierrec/lz4) | Level5 | 7.06ms | 24.94 MB | 39 |
| [LZ4](https://github.com/pierrec/lz4) | Fast | 6.99ms | 24.94 MB | 40 |
| [S2](https://github.com/klauspost/compress) | Better | 7.01ms | 22.95 MB | 34 |
| [S2](https://github.com/klauspost/compress) | Best | 5.05ms | 22.95 MB | 34 |
| [Snappy](https://github.com/golang/snappy) | — | 4.55ms | 8.55 MB | 1 |
| [S2](https://github.com/klauspost/compress) | Default | 8.49ms | 22.95 MB | 34 |
| (none) | — | 694.6µs | 8.55 MB | 1 |

## Notes

- All libraries are pure Go with no CGo dependencies
- Benchmarks run 2 warmup + 10 measured iterations per compressor; median values reported
- Times measured via monotonic clock (time.Since) in timing-only runs
- Memory/allocs measured via runtime.MemStats deltas in separate memory-only runs (after runtime.GC())
- Rows sorted by compressed size ascending within each table
- Ratio = original size / compressed size (higher is better)
- Memory figures are approximate; runtime.MemStats uses process-global counters, so allocations from concurrent goroutines may be included
