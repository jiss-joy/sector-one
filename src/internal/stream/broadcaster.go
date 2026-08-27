package stream

import "sync"

type Broadcaster struct {
	mu   sync.Mutex
	subs map[chan []byte]struct{}
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{subs: make(map[chan []byte]struct{})}
}

func (broadcaster *Broadcaster) Subscribe() chan []byte {
	ch := make(chan []byte, 1)
	broadcaster.mu.Lock()
	broadcaster.subs[ch] = struct{}{}
	broadcaster.mu.Unlock()

	return ch
}

func (broadcaster *Broadcaster) Unsubscribe(ch chan []byte) {
	broadcaster.mu.Lock()
	delete(broadcaster.subs, ch)
	broadcaster.mu.Unlock()
}

func (broadcaster *Broadcaster) Publish(b []byte) {
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()
	for ch := range broadcaster.subs {
		select {
		case ch <- b:
		default:
			select {
			case <-ch:
			default:
			}
			select {
			case ch <- b:
			default:
			}
		}
	}
}

func (broadcaster *Broadcaster) Subscribers() int {
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()
	return len(broadcaster.subs)
}
