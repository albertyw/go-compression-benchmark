package compressors

import (
	"bytes"
	"io"

	"github.com/ulikunitz/xz"
)

func init() {
	Register(Compressor{
		Name:    "XZ/LZMA2",
		Level:   "—",
		Library: "github.com/ulikunitz/xz",
		PureGo:  true,
		Notes:   "Pure Go, fully portable but significantly slower than C implementation",
		Compress: func(data []byte) ([]byte, error) {
			var buf bytes.Buffer
			w, err := xz.NewWriter(&buf)
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
			r, err := xz.NewReader(bytes.NewReader(data))
			if err != nil {
				return nil, err
			}
			return io.ReadAll(r)
		},
	})
}
