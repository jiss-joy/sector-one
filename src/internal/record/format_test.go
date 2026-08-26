package record

import (
	"bytes"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteHeader(&buf); err != nil {
		t.Fatal(err)
	}
	want := Record{Kind: KindCar, Timestamp: 42, Payload: bytes.Repeat([]byte{0xab}, 328)}
	if err := WriteRecord(&buf, want); err != nil {
		t.Fatal(err)
	}
	if err := ReadHeader(&buf); err != nil {
		t.Fatal(err)
	}
	got, err := ReadRecord(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != want.Kind || got.Timestamp != want.Timestamp || !bytes.Equal(got.Payload, want.Payload) {
		t.Fatalf("got %+v", got)
	}
}
