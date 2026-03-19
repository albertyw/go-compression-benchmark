package compressors

import (
	"bytes"
	"compress/gzip"
	"io"
)

func init() {
	levels := []struct {
		name  string
		level int
	}{
		{"BestSpeed", gzip.BestSpeed},
		{"Default", gzip.DefaultCompression},
		{"BestCompression", gzip.BestCompression},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "Gzip stdlib",
			Level:   level.name,
			Library: "compress/gzip",
			PureGo:  true,
			Notes:   "Part of Go stdlib, guaranteed cross-platform support",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w, err := gzip.NewWriterLevel(&buf, level.level)
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
				r, err := gzip.NewReader(bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				defer r.Close()
				return io.ReadAll(r)
			},
		})
	}
}
