package record

import (
	"io"
	"os"
	"time"
)

type Writer struct {
	w io.WriteCloser
}

func NewWriter(path string) (*Writer, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if err := WriteHeader(f); err != nil {
		f.Close()
		return nil, err
	}
	return &Writer{w: f}, nil
}

func (wr *Writer) Write(kind uint8, payload []byte) error {
	return WriteRecord(wr.w, Record{
		Kind:      kind,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	})
}

func (wr *Writer) Close() error { return wr.w.Close() }
