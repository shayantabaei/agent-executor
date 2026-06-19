package execution

type LimitedWriter struct {
	data  []byte
	limit int
}

func NewLimitedWriter(limit int) *LimitedWriter {
	return &LimitedWriter{
		data:  make([]byte, 0, limit),
		limit: limit,
	}
}

func (w *LimitedWriter) Write(p []byte) (int, error) {
	remaining := w.limit - len(w.data)

	if remaining > 0 {
		if len(p) > remaining {
			w.data = append(w.data, p[:remaining]...)
		} else {
			w.data = append(w.data, p...)
		}
	}

	return len(p), nil
}

func (w *LimitedWriter) String() string {
	return string(w.data)
}
