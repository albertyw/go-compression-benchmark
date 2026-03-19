package main

import (
	"errors"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

// --- Roundtrip tests ---

func TestRoundtrip(t *testing.T) {
	testData := []byte("Hello, this is test data for compression roundtrip verification. " +
		"It needs to be long enough to be compressible by most algorithms. " +
		"Repeating patterns help: aaaaabbbbbcccccdddddeeeee aaaaabbbbbcccccdddddeeeee")

	for _, c := range GetAll() {
		t.Run(c.Name+"/"+c.Level, func(t *testing.T) {
			compressed, err := c.Compress(testData)
			if err != nil {
				t.Fatalf("Compress: %v", err)
			}
			decompressed, err := c.Decompress(compressed)
			if err != nil {
				t.Fatalf("Decompress: %v", err)
			}
			if string(decompressed) != string(testData) {
				t.Fatalf("roundtrip mismatch: got %d bytes, want %d bytes", len(decompressed), len(testData))
			}
		})
	}
}

func TestRoundtripEmpty(t *testing.T) {
	testData := []byte{}

	for _, c := range GetAll() {
		t.Run(c.Name+"/"+c.Level, func(t *testing.T) {
			compressed, err := c.Compress(testData)
			if err != nil {
				t.Fatalf("Compress: %v", err)
			}
			decompressed, err := c.Decompress(compressed)
			if err != nil {
				t.Fatalf("Decompress: %v", err)
			}
			if len(decompressed) != 0 {
				t.Fatalf("expected empty output, got %d bytes", len(decompressed))
			}
		})
	}
}

// --- Benchmark runner tests ---

func makeTestCompressor(name string) Compressor {
	return Compressor{
		Name:    name,
		Level:   "test",
		Library: "test/lib",
		PureGo:  true,
		Notes:   "test compressor",
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
	}
}

func TestBenchmarkOne(t *testing.T) {
	c := makeTestCompressor("test")
	data := []byte("test data for benchmarking")

	result, err := benchmarkOne(c, data)
	if err != nil {
		t.Fatalf("benchmarkOne: %v", err)
	}

	if result.OriginalSize != len(data) {
		t.Errorf("OriginalSize = %d, want %d", result.OriginalSize, len(data))
	}
	if result.CompressedSize != len(data) {
		t.Errorf("CompressedSize = %d, want %d (no-op compressor)", result.CompressedSize, len(data))
	}
	if result.Ratio != 1.0 {
		t.Errorf("Ratio = %f, want 1.0", result.Ratio)
	}
	if result.CompressTime == 0 {
		t.Error("CompressTime should be > 0")
	}
	if result.DecompressTime == 0 {
		t.Error("DecompressTime should be > 0")
	}
}

func TestBenchmarkOneCompressError(t *testing.T) {
	c := Compressor{
		Name:  "failing",
		Level: "test",
		Compress: func(data []byte) ([]byte, error) {
			return nil, errors.New("compress failed")
		},
		Decompress: func(data []byte) ([]byte, error) {
			return data, nil
		},
	}
	_, err := benchmarkOne(c, []byte("data"))
	if err == nil {
		t.Fatal("expected error from failing compressor")
	}
	if !strings.Contains(err.Error(), "compress") {
		t.Errorf("error should mention compress, got: %v", err)
	}
}

func TestBenchmarkOneDecompressError(t *testing.T) {
	c := Compressor{
		Name:  "failing-decompress",
		Level: "test",
		Compress: func(data []byte) ([]byte, error) {
			return data, nil
		},
		Decompress: func(data []byte) ([]byte, error) {
			return nil, errors.New("decompress failed")
		},
	}
	_, err := benchmarkOne(c, []byte("data"))
	if err == nil {
		t.Fatal("expected error from failing decompressor")
	}
	if !strings.Contains(err.Error(), "decompress") {
		t.Errorf("error should mention decompress, got: %v", err)
	}
}

func TestBenchmarkOneRoundtripMismatch(t *testing.T) {
	c := Compressor{
		Name:  "mismatch",
		Level: "test",
		Compress: func(data []byte) ([]byte, error) {
			return data, nil
		},
		Decompress: func(data []byte) ([]byte, error) {
			return []byte("wrong"), nil
		},
	}
	_, err := benchmarkOne(c, []byte("data"))
	if err == nil {
		t.Fatal("expected error from roundtrip mismatch")
	}
	if !strings.Contains(err.Error(), "mismatch") {
		t.Errorf("error should mention mismatch, got: %v", err)
	}
}

func TestRunBenchmarks(t *testing.T) {
	saved := registry
	defer func() { registry = saved }()

	registry = nil
	Register(makeTestCompressor("A"))
	Register(makeTestCompressor("B"))

	data := []byte("some test data for RunBenchmarks")
	results := RunBenchmarks(data)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for i := 1; i < len(results); i++ {
		if results[i].CompressedSize < results[i-1].CompressedSize {
			t.Error("results not sorted by compressed size")
		}
	}
}

func TestRunBenchmarksSkipsErrors(t *testing.T) {
	saved := registry
	defer func() { registry = saved }()

	registry = nil
	Register(makeTestCompressor("good"))
	Register(Compressor{
		Name:  "bad",
		Level: "test",
		Compress: func(data []byte) ([]byte, error) {
			return nil, errors.New("fail")
		},
		Decompress: func(data []byte) ([]byte, error) {
			return data, nil
		},
	})

	results := RunBenchmarks([]byte("data"))
	if len(results) != 1 {
		t.Fatalf("expected 1 result (bad skipped), got %d", len(results))
	}
	if results[0].Compressor.Name != "good" {
		t.Errorf("expected 'good' compressor, got %q", results[0].Compressor.Name)
	}
}

// --- Median helper tests ---

func TestMedianDuration(t *testing.T) {
	tests := []struct {
		name string
		in   []time.Duration
		want time.Duration
	}{
		{"odd count", []time.Duration{3, 1, 2}, 2},
		{"even count", []time.Duration{4, 1, 3, 2}, 3},
		{"single", []time.Duration{5}, 5},
		{"already sorted", []time.Duration{1, 2, 3, 4, 5}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := medianDuration(tt.in)
			if got != tt.want {
				t.Errorf("medianDuration(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestMedianUint64(t *testing.T) {
	tests := []struct {
		name string
		in   []uint64
		want uint64
	}{
		{"odd count", []uint64{30, 10, 20}, 20},
		{"even count", []uint64{40, 10, 30, 20}, 30},
		{"single", []uint64{50}, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := medianUint64(tt.in)
			if got != tt.want {
				t.Errorf("medianUint64(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

// --- Format function tests ---

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{1572864, "1.50 MB"},
	}
	for _, tt := range tests {
		got := formatBytes(tt.input)
		if got != tt.want {
			t.Errorf("formatBytes(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatInt(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1,000"},
		{1234567, "1,234,567"},
		{-1234, "-1,234"},
		{-999, "-999"},
		{100, "100"},
	}
	for _, tt := range tests {
		got := formatInt(tt.input)
		if got != tt.want {
			t.Errorf("formatInt(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatRatio(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{1.0, "1.00x"},
		{5.678, "5.68x"},
		{10.333, "10.33x"},
		{0.0, "0.00x"},
	}
	for _, tt := range tests {
		got := formatRatio(tt.input)
		if got != tt.want {
			t.Errorf("formatRatio(%f) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name  string
		input time.Duration
		want  string
	}{
		{"zero", 0, "—"},
		{"microseconds", 500 * time.Microsecond, "500.0µs"},
		{"sub-millisecond", 100 * time.Microsecond, "100.0µs"},
		{"milliseconds", 50 * time.Millisecond, "50.00ms"},
		{"sub-second", 999 * time.Millisecond, "999.00ms"},
		{"seconds", 2 * time.Second, "2.00s"},
		{"seconds fractional", 1500 * time.Millisecond, "1.50s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.input)
			if got != tt.want {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatPureGo(t *testing.T) {
	tests := []struct {
		name string
		c    Compressor
		want string
	}{
		{"empty library", Compressor{Library: ""}, "—"},
		{"pure go", Compressor{Library: "lib", PureGo: true}, "Yes"},
		{"not pure go", Compressor{Library: "lib", PureGo: false}, "No"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatPureGo(tt.c)
			if got != tt.want {
				t.Errorf("formatPureGo() = %q, want %q", got, tt.want)
			}
		})
	}
}

// --- Library link tests ---

func TestLibraryURL(t *testing.T) {
	tests := []struct {
		library string
		want    string
	}{
		{"", ""},
		{"compress/gzip", "https://pkg.go.dev/compress/gzip"},
		{"github.com/klauspost/compress/zstd", "https://github.com/klauspost/compress"},
		{"github.com/klauspost/compress/gzip", "https://github.com/klauspost/compress"},
		{"github.com/klauspost/compress/s2", "https://github.com/klauspost/compress"},
		{"github.com/andybalholm/brotli", "https://github.com/andybalholm/brotli"},
		{"github.com/golang/snappy", "https://github.com/golang/snappy"},
		{"github.com/pierrec/lz4/v4", "https://github.com/pierrec/lz4"},
		{"github.com/ulikunitz/xz", "https://github.com/ulikunitz/xz"},
		{"github.com/klauspost/pgzip", "https://github.com/klauspost/pgzip"},
		{"github.com/klauspost/compress/flate", "https://github.com/klauspost/compress"},
	}
	for _, tt := range tests {
		t.Run(tt.library, func(t *testing.T) {
			got := libraryURL(tt.library)
			if got != tt.want {
				t.Errorf("libraryURL(%q) = %q, want %q", tt.library, got, tt.want)
			}
		})
	}
}

func TestLibraryLink(t *testing.T) {
	tests := []struct {
		name    string
		library string
		want    string
	}{
		{"Zstd", "github.com/klauspost/compress/zstd", "[Zstd](https://github.com/klauspost/compress)"},
		{"(none)", "", "(none)"},
		{"Gzip stdlib", "compress/gzip", "[Gzip stdlib](https://pkg.go.dev/compress/gzip)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := libraryLink(tt.name, tt.library)
			if got != tt.want {
				t.Errorf("libraryLink(%q, %q) = %q, want %q", tt.name, tt.library, got, tt.want)
			}
		})
	}
}

// --- getCPUModel test ---

func TestGetCPUModel(t *testing.T) {
	model := getCPUModel()
	if model == "" {
		t.Error("getCPUModel() returned empty string")
	}
	// On Linux it should return a real model; on other platforms "unknown" is acceptable
}

// --- resultsPath test ---

func TestResultsPath(t *testing.T) {
	path := resultsPath("data.h3")
	if !strings.HasPrefix(path, "results/") {
		t.Errorf("resultsPath() should start with 'results/', got %q", path)
	}
	if !strings.HasSuffix(path, ".md") {
		t.Errorf("resultsPath() should end with '.md', got %q", path)
	}
	// Should contain data file base name (without extension)
	if !strings.Contains(path, "data") {
		t.Errorf("resultsPath() should contain data file name, got %q", path)
	}
	// Should contain arch and OS
	if !strings.Contains(path, runtime.GOARCH) {
		t.Errorf("resultsPath() should contain GOARCH %q, got %q", runtime.GOARCH, path)
	}
	if !strings.Contains(path, runtime.GOOS) {
		t.Errorf("resultsPath() should contain GOOS %q, got %q", runtime.GOOS, path)
	}
	// Should contain a date pattern
	if !strings.Contains(path, "-") {
		t.Errorf("resultsPath() should contain date separators, got %q", path)
	}
}

// --- WriteResults tests ---

func TestWriteResults(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	section := ResultSection{
		FileName: "test.dat",
		FileSize: 1000,
		Results: []BenchmarkResult{
			{
				Compressor: Compressor{
					Name:    "TestLib",
					Level:   "fast",
					Library: "github.com/example/testlib",
					PureGo:  true,
					Notes:   "test notes",
				},
				OriginalSize:     1000,
				CompressedSize:   500,
				Ratio:            2.0,
				CompressTime:     10 * time.Millisecond,
				DecompressTime:   5 * time.Millisecond,
				CompressMem:      1024,
				CompressAllocs:   10,
				DecompressMem:    512,
				DecompressAllocs: 5,
			},
			{
				Compressor: Compressor{
					Name:    "NoLib",
					Level:   "—",
					Library: "",
					PureGo:  false,
					Notes:   "baseline",
				},
				OriginalSize:   1000,
				CompressedSize: 1000,
				Ratio:          1.0,
			},
		},
	}

	path, err := WriteResults(section)
	if err != nil {
		t.Fatalf("WriteResults: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	s := string(content)

	// Check environment section
	if !strings.Contains(s, "## Environment") {
		t.Error("missing Environment section")
	}
	if !strings.Contains(s, "**CPU:**") {
		t.Error("missing CPU info")
	}
	if !strings.Contains(s, "**Architecture:**") {
		t.Error("missing Architecture info")
	}
	if !strings.Contains(s, "**OS:**") {
		t.Error("missing OS info")
	}
	if !strings.Contains(s, "**Go:**") {
		t.Error("missing Go version")
	}

	// Check structure
	if !strings.Contains(s, "# Compression Benchmark Results") {
		t.Error("missing title")
	}
	if !strings.Contains(s, "## test.dat (1,000 bytes)") {
		t.Error("missing section header")
	}

	// Check separate compression and decompression tables
	if !strings.Contains(s, "### Compression") {
		t.Error("missing Compression subsection")
	}
	if !strings.Contains(s, "### Decompression") {
		t.Error("missing Decompression subsection")
	}

	// Check library link
	if !strings.Contains(s, "[TestLib](https://github.com/example/testlib)") {
		t.Error("missing TestLib github link")
	}

	// NoLib (no library) should not be linked
	if !strings.Contains(s, "| NoLib |") {
		t.Error("missing NoLib row")
	}

	// Check PureGo column values
	if !strings.Contains(s, "| Yes |") {
		t.Error("missing Pure Go Yes")
	}
	if !strings.Contains(s, "| — | baseline |") {
		t.Error("missing baseline with — for PureGo")
	}

	// Notes section
	if !strings.Contains(s, "## Notes") {
		t.Error("missing notes section")
	}
	if !strings.Contains(s, "Ratio = original size / compressed size") {
		t.Error("missing ratio explanation")
	}
}

func TestWriteResultsPureGoNo(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	section := ResultSection{
		FileName: "test.dat",
		FileSize: 100,
		Results: []BenchmarkResult{
			{
				Compressor: Compressor{
					Name:    "CgoLib",
					Level:   "default",
					Library: "cgo/lib",
					PureGo:  false,
				},
			},
		},
	}

	path, err := WriteResults(section)
	if err != nil {
		t.Fatalf("WriteResults: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(content), "| No |") {
		t.Error("expected PureGo = No for CGo library")
	}
}
