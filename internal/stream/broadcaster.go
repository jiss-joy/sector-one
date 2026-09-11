package stream

import "sync"

type Broadcaster struct {
	mutex       sync.Mutex
	subscribers map[chan []byte]struct{}
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		subscribers: make(map[chan []byte]struct{}),
	}
}

func (broadcaster *Broadcaster) Subscribe() chan []byte {
	ch := make(chan []byte, 1)
	broadcaster.mutex.Lock()
	broadcaster.subscribers[ch] = struct{}{}
	broadcaster.mutex.Unlock()

	return ch
}

func (broadcaster *Broadcaster) Unsubscribe(ch chan []byte) {
	broadcaster.mutex.Lock()
	delete(broadcaster.subscribers, ch)
	broadcaster.mutex.Unlock()
}

func (broadcaster *Broadcaster) Subscribers() int {
	broadcaster.mutex.Lock()
	defer broadcaster.mutex.Unlock()
	return len(broadcaster.subscribers)
}

func (broadcaster *Broadcaster) Publish(bytes []byte) {
	broadcaster.mutex.Lock()
	defer broadcaster.mutex.Unlock()

	for subscriber := range broadcaster.subscribers {
		select {
		case subscriber <- bytes:
		default:
			select {
			case <-subscriber:
			default:
			}
			select {
			case subscriber <- bytes:
			default:
			}
		}
	}
}
