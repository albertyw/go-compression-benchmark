package main

import (
	"bytes"
	"io"

	kgzip "github.com/klauspost/compress/gzip"
)

func init() {
	levels := []struct {
		name  string
		level int
	}{
		{"BestSpeed", kgzip.BestSpeed},
		{"Default", kgzip.DefaultCompression},
		{"BestCompression", kgzip.BestCompression},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "Gzip klauspost",
			Level:   level.name,
			Library: "github.com/klauspost/compress/gzip",
			PureGo:  true,
			Notes:   "Drop-in stdlib replacement, assembly optimizations on amd64",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w, err := kgzip.NewWriterLevel(&buf, level.level)
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
				r, err := kgzip.NewReader(bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				defer r.Close()
				return io.ReadAll(r)
			},
		})
	}
}
