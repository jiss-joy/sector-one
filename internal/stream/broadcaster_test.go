package stream

import (
	"bytes"
	"testing"
	"time"
)

func TestBroadcasterDroppedLatest(t *testing.T) {
	h := NewBroadcaster()
	ch := h.Subscribe()
	defer h.Unsubscribe(ch)

	h.Publish([]byte("a"))
	h.Publish([]byte("b"))

	select {
	case got := <-ch:
		if !bytes.Equal(got, []byte("b")) {
			t.Fatalf("got %q, want b", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestBroadcasterUnsubscribe(t *testing.T) {
	h := NewBroadcaster()
	ch := h.Subscribe()
	if h.Subscribers() != 1 {
		t.Fatalf("subs=%d", h.Subscribers())
	}
	h.Unsubscribe(ch)
	if h.Subscribers() != 0 {
		t.Fatalf("subs=%d after unsub", h.Subscribers())
	}
}
