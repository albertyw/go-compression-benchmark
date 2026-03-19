package compressors

import (
	"bytes"
	"io"

	"github.com/pierrec/lz4/v4"
)

func init() {
	levels := []struct {
		name  string
		level lz4.CompressionLevel
	}{
		{"Fast", lz4.Fast},
		{"Level5", lz4.Level5},
		{"Level9", lz4.Level9},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "LZ4",
			Level:   level.name,
			Library: "github.com/pierrec/lz4/v4",
			PureGo:  true,
			Notes:   "Pure Go, fully portable",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w := lz4.NewWriter(&buf)
				if err := w.Apply(lz4.CompressionLevelOption(level.level)); err != nil {
					return nil, err
				}
				if _, err := w.Write(data); err != nil {
					return nil, err
				}
				if err := w.Close(); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				r := lz4.NewReader(bytes.NewReader(data))
				return io.ReadAll(r)
			},
		})
	}
}
