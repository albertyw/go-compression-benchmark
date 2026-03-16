package main

func init() {
	Register(Compressor{
		Name:    "(none)",
		Level:   "—",
		Library: "",
		PureGo:  false,
		Notes:   "Baseline — no compression",
		Compress: func(data []byte) ([]byte, error) {
			out := make([]byte, len(data))
			copy(out, data)
			return out, nil
		},
		Decompress: func(data []byte) ([]byte, error) {
			out := make([]byte, len(data))
			copy(out, data)
			return out, nil
		},
	})
}
