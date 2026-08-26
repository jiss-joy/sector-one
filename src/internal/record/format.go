package record

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	Version uint32 = 1
	KindHello  uint8  = 1
	KindCar    uint8  = 2
	maxPayload = 64 * 1024
)

var magic = [8]byte{'S', '1', 'R', 'E', 'C', 0, 0, 1}

type Record struct {
	Kind      uint8
	Timestamp int64
	Payload   []byte
}

func WriteHeader(w io.Writer) error {
	var hdr [16]byte
	copy(hdr[0:8], magic[:])
	binary.LittleEndian.PutUint32(hdr[8:12], Version)
	_, err := w.Write(hdr[:])
	return err
}

func ReadHeader(r io.Reader) error {
	var hdr [16]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return err
	}
	if string(hdr[0:8]) != string(magic[:]) {
		return fmt.Errorf("record: bad magic")
	}
	if binary.LittleEndian.Uint32(hdr[8:12]) != Version {
		return fmt.Errorf("record: unsupported version")
	}
	return nil
}

func WriteRecord(w io.Writer, rec Record) error {
	if len(rec.Payload) > maxPayload {
		return fmt.Errorf("record: payload too large")
	}
	var prefix [13]byte
	prefix[0] = rec.Kind
	binary.LittleEndian.PutUint64(prefix[1:9], uint64(rec.Timestamp))
	binary.LittleEndian.PutUint32(prefix[9:13], uint32(len(rec.Payload)))
	if _, err := w.Write(prefix[:]); err != nil {
		return err
	}
	_, err := w.Write(rec.Payload)
	return err
}

func ReadRecord(r io.Reader) (Record, error) {
	var prefix [13]byte
	if _, err := io.ReadFull(r, prefix[:]); err != nil {
		return Record{}, err
	}
	n := binary.LittleEndian.Uint32(prefix[9:13])
	if n > maxPayload {
		return Record{}, fmt.Errorf("record: payload %d too large", n)
	}
	payload := make([]byte, n)
	if _, err := io.ReadFull(r, payload); err != nil {
		return Record{}, err
	}
	return Record{
		Kind:      prefix[0],
		Timestamp: int64(binary.LittleEndian.Uint64(prefix[1:9])),
		Payload:   payload,
	}, nil
}
