package main

import (
	"github.com/golang/snappy"
)

func init() {
	Register(Compressor{
		Name:    "Snappy",
		Level:   "—",
		Library: "github.com/golang/snappy",
		PureGo:  true,
		Notes:   "Pure Go, assembly on amd64",
		Compress: func(data []byte) ([]byte, error) {
			return snappy.Encode(nil, data), nil
		},
		Decompress: func(data []byte) ([]byte, error) {
			return snappy.Decode(nil, data)
		},
	})
}
