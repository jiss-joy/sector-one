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
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	err = WriteHeader(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	return &Writer{w: file}, nil
}

func (writer *Writer) Write(kind uint8, payload []byte) error {
	return WriteRecord(writer.w, Record{
		Kind:      kind,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	})
}

func (writer *Writer) Close() error { return writer.w.Close() }
