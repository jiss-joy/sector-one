package record

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestReaderFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "clip.bin")
	wr, err := NewWriter(path)
	if err != nil {
		t.Fatal(err)
	}
	hello := bytes.Repeat([]byte{1}, 408)
	car := bytes.Repeat([]byte{'a'}, 328)
	if err := wr.Write(KindHello, hello); err != nil {
		t.Fatal(err)
	}
	if err := wr.Write(KindCar, car); err != nil {
		t.Fatal(err)
	}
	if err := wr.Close(); err != nil {
		t.Fatal(err)
	}

	rd, err := NewReader(path)
	if err != nil {
		t.Fatal(err)
	}
	defer rd.Close()

	r1, err := rd.Next()
	if err != nil {
		t.Fatal(err)
	}
	if r1.Kind != KindHello || !bytes.Equal(r1.Payload, hello) {
		t.Fatalf("hello: %+v", r1)
	}
	r2, err := rd.Next()
	if err != nil {
		t.Fatal(err)
	}
	if r2.Kind != KindCar || !bytes.Equal(r2.Payload, car) {
		t.Fatalf("car: %+v", r2)
	}
	if _, err := rd.Next(); err != io.EOF {
		t.Fatalf("want EOF, got %v", err)
	}
}

func TestReaderBadMagic(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.bin")
	if err := os.WriteFile(path, []byte("not a recording"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := NewReader(path); err == nil {
		t.Fatal("expected bad magic")
	}
}
