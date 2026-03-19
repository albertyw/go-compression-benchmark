package compressors

import (
	"bytes"
	"io"

	kflate "github.com/klauspost/compress/flate"
)

func init() {
	levels := []struct {
		name  string
		level int
	}{
		{"BestSpeed", kflate.BestSpeed},
		{"Default", kflate.DefaultCompression},
		{"BestCompression", kflate.BestCompression},
	}
	for _, l := range levels {
		level := l
		Register(Compressor{
			Name:    "Deflate klauspost",
			Level:   level.name,
			Library: "github.com/klauspost/compress/flate",
			PureGo:  true,
			Notes:   "Drop-in stdlib replacement, assembly optimizations on amd64",
			Compress: func(data []byte) ([]byte, error) {
				var buf bytes.Buffer
				w, err := kflate.NewWriter(&buf, level.level)
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
				r := kflate.NewReader(bytes.NewReader(data))
				defer r.Close()
				return io.ReadAll(r)
			},
		})
	}
}
