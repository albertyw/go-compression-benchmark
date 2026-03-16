# Compression Benchmark Results

## data.h3 (8,966,666 bytes)

| Library | Level | Compressed Size | Ratio | Compress Time | Decompress Time | Compress Mem | Compress Allocs | Decompress Mem | Decompress Allocs | Pure Go | Notes |
|---------|-------|----------------|-------|--------------|----------------|-------------|----------------|---------------|-------------------|---------|-------|
| XZ/LZMA2 | — | 847.99 KB | 10.33x | 448.82ms | 97.21ms | 64.76 MB | 932,861 | 43.18 MB | 931,689 | Yes | Pure Go, fully portable but significantly slower than C implementation |
| Brotli | Best (11) | 1.10 MB | 7.79x | 10.10s | 24.53ms | 802.51 MB | 373 | 26.30 MB | 54 | Yes | Pure Go, fully portable |
| Brotli | Default (6) | 1.20 MB | 7.15x | 272.45ms | 23.77ms | 338.96 MB | 131 | 25.08 MB | 53 | Yes | Pure Go, fully portable |
| Zstd | SpeedBestCompression | 1.28 MB | 6.67x | 317.07ms | 13.94ms | 23.22 MB | 63 | 8.56 MB | 7 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| Zstd | SpeedBetterCompression | 1.47 MB | 5.83x | 61.83ms | 15.37ms | 26.16 MB | 66 | 8.56 MB | 7 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| Zstd | SpeedFastest | 1.48 MB | 5.78x | 34.01ms | 13.16ms | 16.55 MB | 63 | 8.56 MB | 7 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| Zstd | SpeedDefault | 1.52 MB | 5.64x | 37.50ms | 14.96ms | 26.16 MB | 66 | 8.56 MB | 7 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| Brotli | Fast (1) | 1.60 MB | 5.36x | 39.28ms | 38.82ms | 14.90 MB | 14 | 29.53 MB | 662 | Yes | Pure Go, fully portable |
| Gzip stdlib | Default | 1.60 MB | 5.36x | 172.25ms | 21.59ms | 4.78 MB | 33 | 21.68 MB | 3,624 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| Gzip stdlib | BestCompression | 1.60 MB | 5.35x | 284.20ms | 20.82ms | 4.78 MB | 33 | 21.68 MB | 3,559 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| pgzip | BestCompression | 1.60 MB | 5.35x | 46.47ms | 15.78ms | 45.64 MB | 302 | 24.99 MB | 196 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| Deflate klauspost | BestCompression | 1.60 MB | 5.35x | 285.32ms | 13.18ms | 5.09 MB | 29 | 20.99 MB | 170 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| Gzip klauspost | BestCompression | 1.60 MB | 5.35x | 286.62ms | 13.17ms | 5.09 MB | 31 | 20.99 MB | 172 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| Deflate klauspost | Default | 1.62 MB | 5.26x | 33.65ms | 13.91ms | 5.03 MB | 28 | 21.02 MB | 472 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| Gzip klauspost | Default | 1.62 MB | 5.26x | 33.85ms | 13.63ms | 5.03 MB | 30 | 21.02 MB | 472 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| pgzip | Default | 1.62 MB | 5.26x | 8.14ms | 13.74ms | 42.06 MB | 290 | 25.02 MB | 491 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| Deflate klauspost | BestSpeed | 1.66 MB | 5.16x | 18.85ms | 13.26ms | 4.78 MB | 28 | 21.03 MB | 537 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| Gzip klauspost | BestSpeed | 1.66 MB | 5.16x | 19.02ms | 13.46ms | 4.78 MB | 30 | 21.03 MB | 538 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| pgzip | BestSpeed | 1.66 MB | 5.16x | 5.63ms | 13.81ms | 40.04 MB | 294 | 25.03 MB | 570 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| Gzip stdlib | BestSpeed | 1.91 MB | 4.47x | 33.30ms | 22.57ms | 5.15 MB | 35 | 21.65 MB | 4,284 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| LZ4 | Level9 | 3.52 MB | 2.43x | 105.18ms | 6.10ms | 20.61 MB | 22 | 24.94 MB | 39 | Yes | Pure Go, fully portable |
| LZ4 | Level5 | 3.52 MB | 2.43x | 105.75ms | 8.36ms | 20.61 MB | 21 | 24.94 MB | 40 | Yes | Pure Go, fully portable |
| LZ4 | Fast | 3.55 MB | 2.41x | 15.35ms | 5.96ms | 20.09 MB | 23 | 24.94 MB | 40 | Yes | Pure Go, fully portable |
| S2 | Better | 3.56 MB | 2.40x | 5.88ms | 6.34ms | 28.32 MB | 144 | 22.95 MB | 34 | Yes | Snappy-compatible decoding; S2-compressed data cannot be decoded by Snappy decoders |
| S2 | Best | 3.57 MB | 2.40x | 53.26ms | 5.64ms | 72.85 MB | 155 | 22.95 MB | 34 | Yes | Snappy-compatible decoding; S2-compressed data cannot be decoded by Snappy decoders |
| Snappy | — | 3.91 MB | 2.19x | 8.09ms | 6.18ms | 9.98 MB | 1 | 8.55 MB | 1 | Yes | Pure Go, assembly on amd64 |
| S2 | Default | 4.14 MB | 2.07x | 6.10ms | 6.51ms | 31.62 MB | 145 | 22.95 MB | 34 | Yes | Snappy-compatible decoding; S2-compressed data cannot be decoded by Snappy decoders |
| (none) | — | 8.55 MB | 1.00x | 336.2µs | 578.4µs | 8.55 MB | 1 | 8.55 MB | 1 | — | Baseline — no compression |

## data_mock.h3 (1,253 bytes)

| Library | Level | Compressed Size | Ratio | Compress Time | Decompress Time | Compress Mem | Compress Allocs | Decompress Mem | Decompress Allocs | Pure Go | Notes |
|---------|-------|----------------|-------|--------------|----------------|-------------|----------------|---------------|-------------------|---------|-------|
| Brotli | Best (11) | 142 B | 8.82x | 1.12ms | 11.4µs | 33.01 MB | 61 | 73.70 KB | 19 | Yes | Pure Go, fully portable |
| Brotli | Default (6) | 156 B | 8.03x | 105.7µs | 11.7µs | 2.09 MB | 40 | 65.76 KB | 19 | Yes | Pure Go, fully portable |
| Zstd | SpeedFastest | 191 B | 6.56x | 37.7µs | 5.2µs | 1.25 KB | 1 | 7.70 KB | 7 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| Zstd | SpeedDefault | 191 B | 6.56x | 15.0µs | 2.7µs | 1.25 KB | 1 | 7.70 KB | 7 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| Zstd | SpeedBetterCompression | 191 B | 6.56x | 24.1µs | 3.0µs | 1.25 KB | 1 | 7.70 KB | 7 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| Zstd | SpeedBestCompression | 194 B | 6.46x | 58.5µs | 3.3µs | 1.25 KB | 1 | 7.70 KB | 7 | Yes | Pure Go, assembly on amd64, good cross-platform support |
| Deflate klauspost | BestCompression | 232 B | 5.40x | 58.4µs | 5.6µs | 1.09 MB | 16 | 40.08 KB | 12 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| Deflate klauspost | Default | 242 B | 5.18x | 39.9µs | 4.8µs | 1.03 MB | 15 | 40.08 KB | 12 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| Brotli | Fast (1) | 242 B | 5.18x | 9.8µs | 11.4µs | 44.16 KB | 12 | 65.76 KB | 19 | Yes | Pure Go, fully portable |
| Deflate klauspost | BestSpeed | 248 B | 5.05x | 34.2µs | 5.3µs | 794.92 KB | 15 | 40.08 KB | 12 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| Gzip klauspost | BestCompression | 250 B | 5.01x | 58.8µs | 5.9µs | 1.09 MB | 18 | 40.77 KB | 13 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| Gzip stdlib | Default | 252 B | 4.97x | 43.3µs | 5.8µs | 795.59 KB | 21 | 43.20 KB | 11 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| Gzip stdlib | BestCompression | 252 B | 4.97x | 48.5µs | 6.3µs | 795.59 KB | 21 | 43.20 KB | 11 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| pgzip | BestCompression | 255 B | 4.91x | 221.0µs | 188.0µs | 4.29 MB | 46 | 4.04 MB | 29 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| Gzip klauspost | Default | 260 B | 4.82x | 41.7µs | 5.6µs | 1.03 MB | 18 | 40.77 KB | 13 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| Gzip stdlib | BestSpeed | 263 B | 4.76x | 42.9µs | 6.1µs | 1.15 MB | 23 | 43.20 KB | 11 | Yes | Part of Go stdlib, guaranteed cross-platform support |
| pgzip | Default | 265 B | 4.73x | 206.0µs | 185.3µs | 4.22 MB | 45 | 4.04 MB | 29 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| Gzip klauspost | BestSpeed | 266 B | 4.71x | 29.7µs | 5.4µs | 795.64 KB | 18 | 40.77 KB | 13 | Yes | Drop-in stdlib replacement, assembly optimizations on amd64 |
| pgzip | BestSpeed | 271 B | 4.62x | 187.9µs | 179.2µs | 3.97 MB | 45 | 4.04 MB | 29 | Yes | Parallel gzip — performance scales with CPU cores, results vary by machine |
| XZ/LZMA2 | — | 272 B | 4.61x | 1.44ms | 32.2µs | 48.13 MB | 414 | 8.04 MB | 261 | Yes | Pure Go, fully portable but significantly slower than C implementation |
| Snappy | — | 525 B | 2.39x | 1.0µs | 0.6µs | 1.50 KB | 1 | 1.25 KB | 1 | Yes | Pure Go, assembly on amd64 |
| LZ4 | Level9 | 556 B | 2.25x | 10.9µs | 1.7µs | 5.02 MB | 14 | 4.01 MB | 14 | Yes | Pure Go, fully portable |
| LZ4 | Level5 | 556 B | 2.25x | 10.8µs | 1.8µs | 5.02 MB | 14 | 4.01 MB | 14 | Yes | Pure Go, fully portable |
| S2 | Best | 564 B | 2.22x | 303.8µs | 35.4µs | 6.51 MB | 18 | 1.01 MB | 9 | Yes | Snappy-compatible decoding; S2-compressed data cannot be decoded by Snappy decoders |
| S2 | Better | 565 B | 2.22x | 81.1µs | 34.3µs | 2.03 MB | 19 | 1.01 MB | 9 | Yes | Snappy-compatible decoding; S2-compressed data cannot be decoded by Snappy decoders |
| LZ4 | Fast | 568 B | 2.21x | 1.7µs | 1.5µs | 4.14 MB | 14 | 4.01 MB | 14 | Yes | Pure Go, fully portable |
| S2 | Default | 581 B | 2.16x | 83.6µs | 35.5µs | 2.02 MB | 19 | 1.01 MB | 9 | Yes | Snappy-compatible decoding; S2-compressed data cannot be decoded by Snappy decoders |
| (none) | — | 1.22 KB | 1.00x | 0.1µs | 0.0µs | 1.25 KB | 1 | 1.25 KB | 1 | — | Baseline — no compression |

## Notes

- All libraries are pure Go with no CGo dependencies
- Benchmarks run 2 warmup + 10 measured iterations per compressor; median values reported
- Times measured via monotonic clock (time.Since) in timing-only runs
- Memory/allocs measured via runtime.MemStats deltas in separate memory-only runs (after runtime.GC())
- Rows sorted by compressed size ascending within each table
