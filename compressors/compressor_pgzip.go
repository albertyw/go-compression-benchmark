package compressors

import (
	"bytes"
	"io"

	"github.com/klauspost/pgzip"
)

func init() {
	levels := []struct {
		name  string
		level int
	}{
		{"BestSpeed", pgzip.BestSpeed},
		{"Default", pgzip.DefaultCompression},
		{"BestCompression", pgzip.BestCompression},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "pgzip",
			Level:   level.name,
			Library: "github.com/klauspost/pgzip",
			PureGo:  true,
			Notes:   "Parallel gzip — performance scales with CPU cores, results vary by machine",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w, err := pgzip.NewWriterLevel(&buf, level.level)
				if err != nil {
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
				r, err := pgzip.NewReader(bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				defer r.Close()
				return io.ReadAll(r)
			},
		})
	}
}
