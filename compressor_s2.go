package main

import (
	"bytes"
	"io"

	"github.com/klauspost/compress/s2"
)

func init() {
	type s2Level struct {
		name    string
		options []s2.WriterOption
	}
	levels := []s2Level{
		{"Default", nil},
		{"Better", []s2.WriterOption{s2.WriterBetterCompression()}},
		{"Best", []s2.WriterOption{s2.WriterBestCompression()}},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "S2",
			Level:   level.name,
			Library: "github.com/klauspost/compress/s2",
			PureGo:  true,
			Notes:   "S2 can decode Snappy data, but Snappy decoders cannot decode S2 output",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w := s2.NewWriter(&buf, level.options...)
				if _, err := w.Write(data); err != nil {
					return nil, err
				}
				if err := w.Close(); err != nil {
					return nil, err
				}
				return buf.Bytes(), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				r := s2.NewReader(bytes.NewReader(data))
				return io.ReadAll(r)
			},
		})
	}
}
