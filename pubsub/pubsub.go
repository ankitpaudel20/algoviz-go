package pubsub

import (
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
)

// PubSub struct holds the information about individual topics.

type PubSub[K comparable, V interface{}] struct {
	// this funny type is to store channels as set so
	// that they can be found faster for deletion
	topics map[K]map[chan V]uuid.UUID
	mu     sync.RWMutex
}

// NewPubSub initializes the PubSub instance.
func NewPubSub[K comparable, V interface{}]() *PubSub[K, V] {
	return &PubSub[K, V]{
		topics: make(map[K]map[chan V]uuid.UUID),
	}
}

// Subscribe returns a read-only channel for the given topic.
func (ps *PubSub[K, V]) Subscribe(topic K, chan_capacity int) (chan V, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	new_chan := make(chan V, chan_capacity)
	if chs := ps.topics[topic]; chs == nil {
		ps.topics[topic] = make(map[chan V]uuid.UUID)
	}
	uid := uuid.New()
	ps.topics[topic][new_chan] = uid
	log.Printf("Subscribed to topic: %+v | channel_id: %+v | number of channels: %d", topic, uid, len(ps.topics[topic]))

	return new_chan, nil
}

func (ps *PubSub[K, V]) Unsubscribe(topic K, commChan chan V) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	chs, ok := ps.topics[topic]
	if !ok {
		log.Printf("no such topic exists")
		return
	}
	uid, ok := ps.topics[topic][commChan]
	if !ok {
		log.Printf("no such channel")
		return
	}
	delete(chs, commChan)
	log.Printf("unsubscribed from a topic: %+v | channel_id: %+v | number of channels: %d", topic, uid, len(ps.topics[topic]))
}

// Publish sends data to the given topic.
func (ps *PubSub[K, V]) Publish(topic K, data V) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	chs, ok := ps.topics[topic]
	if !ok {
		return errors.New("topic does not exist")
	}

	for ch, uid := range chs {
		log.Printf("sending data %+v to channel %+v", data, uid)
		ch <- data
	}

	return nil
}
