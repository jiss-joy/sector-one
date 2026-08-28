package stream

import "sync"

type Broadcaster struct {
	mu   sync.Mutex
	subscribers map[chan []byte]struct{}
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{subscribers: make(map[chan []byte]struct{})}
}

func (broadcaster *Broadcaster) Subscribe() chan []byte {
	ch := make(chan []byte, 1)
	broadcaster.mu.Lock()
	broadcaster.subscribers[ch] = struct{}{}
	broadcaster.mu.Unlock()

	return ch
}

func (broadcaster *Broadcaster) Unsubscribe(ch chan []byte) {
	broadcaster.mu.Lock()
	delete(broadcaster.subscribers, ch)
	broadcaster.mu.Unlock()
}

func (broadcaster *Broadcaster) Publish(b []byte) {
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()
	for sub := range broadcaster.subscribers {
		select {
		case sub <- b:
		default:
			select {
			case <-sub:
			default:
			}
			select {
			case sub <- b:
			default:
			}
		}
	}
}

func (broadcaster *Broadcaster) Subscribers() int {
	broadcaster.mu.Lock()
	defer broadcaster.mu.Unlock()
	return len(broadcaster.subscribers)
}
