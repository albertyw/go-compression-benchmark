package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/andybalholm/brotli"
)

func init() {
	levels := []struct {
		name  string
		level int
	}{
		{"Fast", 1},
		{"Default", 6},
		{"Best", 11},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "Brotli",
			Level:   fmt.Sprintf("%s (%d)", level.name, level.level),
			Library: "github.com/andybalholm/brotli",
			PureGo:  true,
			Notes:   "Pure Go, fully portable",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w := brotli.NewWriterLevel(&buf, level.level)
				if _, err := w.Write(data); err != nil {
					return nil, err
				}
				if err := w.Close(); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				r := brotli.NewReader(bytes.NewReader(data))
				return io.ReadAll(r)
			},
		})
	}
}
