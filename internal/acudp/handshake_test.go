package acudp

import (
	"encoding/binary"
	"testing"
	"unicode/utf16"
)

func TestEncodeHandshake(t *testing.T) {
	got := EncodeHandshake(OPSubscribeUpdate)
	if len(got) != 12 {
		t.Fatalf("len=%d, want 12", len(got))
	}
	if binary.LittleEndian.Uint32(got[0:4]) != 0 {
		t.Fatalf("identifier=%d, want 0", binary.LittleEndian.Uint32(got[0:4]))
	}
	if binary.LittleEndian.Uint32(got[4:8]) != 1 {
		t.Fatalf("version=%d, want 1", binary.LittleEndian.Uint32(got[4:8]))
	}
	if binary.LittleEndian.Uint32(got[8:12]) != 1 {
		t.Fatalf("operation=%d, want 1", binary.LittleEndian.Uint32(got[8:12]))
	}
}

func TestParseHandshakeResponse408(t *testing.T) {
	buf := make([]byte, 408)
	n := (408 - 8) / 4 // 50 wchars, 100 bytes each
	putACString(buf[0:n], "ks_porsche%")
	putACString(buf[n:2*n], "jiss%")
	binary.LittleEndian.PutUint32(buf[2*n:2*n+4], 4242)
	binary.LittleEndian.PutUint32(buf[2*n+4:2*n+8], 1)
	putACString(buf[2*n+8:3*n+8], "monza%")
	putACString(buf[3*n+8:4*n+8], "")

	got, err := ParseHandshakeResponse(buf)
	if err != nil {
		t.Fatal(err)
	}
	if got.CarName != "ks_porsche" || got.DriverName != "jiss" || got.TrackName != "monza" {
		t.Fatalf("got %+v", got)
	}
}

func putACString(dst []byte, s string) {
	u16s := utf16.Encode([]rune(s))
	for i, u := range u16s {
		if i*2+1 >= len(dst) {
			return
		}
		binary.LittleEndian.PutUint16(dst[i*2:], u)
	}
}
