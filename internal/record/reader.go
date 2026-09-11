package record

import (
	"io"
	"os"
)

type Reader struct {
	r io.ReadCloser
}

func NewReader(path string) (*Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	err = ReadHeader(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	return &Reader{r: file}, nil
}

func (reader *Reader) Next() (Record, error) {
	return ReadRecord(reader.r)
}

func (reader *Reader) Close() error {
	return reader.r.Close()
}
