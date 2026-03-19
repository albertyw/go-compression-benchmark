package compressors

// Compressor defines a compression algorithm with its metadata and functions.
type Compressor struct {
	Name       string                       // display name, e.g. "Zstd"
	Level      string                       // compression level, e.g. "SpeedFastest"
	Library    string                       // import path, e.g. "github.com/klauspost/compress/zstd"
	PureGo     bool                         // true if no CGo dependency
	Notes      string                       // portability/edge case notes
	Compress   func([]byte) ([]byte, error) // compress input bytes
	Decompress func([]byte) ([]byte, error) // decompress back to original
}

var Registry []Compressor

// Register adds a compressor to the global registry.
func Register(c Compressor) {
	Registry = append(Registry, c)
}

// GetAll returns all registered compressors.
func GetAll() []Compressor {
	return Registry
}
