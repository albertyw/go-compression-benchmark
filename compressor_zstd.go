package main

import (
	"github.com/klauspost/compress/zstd"
)

func init() {
	// Shared decoder — zstd decoding is level-independent
	dec, err := zstd.NewReader(nil)
	if err != nil {
		panic("zstd: failed to create decoder: " + err.Error())
	}

	levels := []struct {
		name  string
		level zstd.EncoderLevel
	}{
		{"SpeedFastest", zstd.SpeedFastest},
		{"SpeedDefault", zstd.SpeedDefault},
		{"SpeedBetterCompression", zstd.SpeedBetterCompression},
		{"SpeedBestCompression", zstd.SpeedBestCompression},
	}
	for _, l := range levels {
		level := l
		enc, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(level.level))
		if err != nil {
			panic("zstd: failed to create encoder: " + err.Error())
		}
		Register(Compressor{
			Name:    "Zstd",
			Level:   level.name,
			Library: "github.com/klauspost/compress/zstd",
			PureGo:  true,
			Notes:   "Pure Go, assembly on amd64, good cross-platform support",
			Compress: func(data []byte) ([]byte, error) {
				return enc.EncodeAll(data, nil), nil
			},
			Decompress: func(data []byte) ([]byte, error) {
				return dec.DecodeAll(data, nil)
			},
		})
	}
}
