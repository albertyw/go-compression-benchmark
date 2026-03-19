package main

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"time"

	"github.com/albertyw/go-compression-benchmark/compressors"
)

const (
	warmupIterations  = 2
	measureIterations = 10
)

// BenchmarkResult holds the median metrics for one compressor on one data file.
type BenchmarkResult struct {
	Compressor       compressors.Compressor
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
	for _, c := range compressors.GetAll() {
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

func benchmarkOne(c compressors.Compressor, data []byte) (BenchmarkResult, error) {
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
