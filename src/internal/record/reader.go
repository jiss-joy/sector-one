package record

import (
	"io"
	"os"
)

type Reader struct {
	r io.ReadCloser
}

func NewReader(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if err := ReadHeader(f); err != nil {
		f.Close()
		return nil, err
	}
	return &Reader{r: f}, nil
}

func (rd *Reader) Next() (Record, error) {
	return ReadRecord(rd.r)
}

func (rd *Reader) Close() error {
	return rd.r.Close()
}
