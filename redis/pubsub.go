package redis

import (
	"github.com/go-redis/redis"
	"github.com/jacygao/broadcaster"
	"sync"
)

type redisBroadcaster struct {
	sync.Mutex
	client    *redis.Client
	pubSub    *redis.PubSub
	callback  map[string]broadcaster.Callback
	bufLen    int
	closeChan chan struct{}
}

// NewRedisBroadcaster returns an implementation of the Broadcaster interface backed by Redis PubSub.
func NewRedisBroadcaster(cli *redis.Client, bufLen int) (*redisBroadcaster, error) {
	ps := cli.Subscribe()
	rb := &redisBroadcaster{
		client:    cli,
		pubSub:    ps,
		callback:  make(map[string]broadcaster.Callback, bufLen),
		bufLen:    bufLen,
		closeChan: make(chan struct{}, 1),
	}

	go rb.run()

	return rb, nil
}

func (rb *redisBroadcaster) run() {
	for {
		select {
		case d := <-rb.pubSub.Channel():
			if d != nil {
				f, ok := rb.callback[d.Channel]
				if ok {
					f(d.Payload)
				}
			}
		case <-rb.closeChan:
			return
		}
	}
}

// Publish publishes a new notification to a given redis channel.
func (rb *redisBroadcaster) Publish(channel string, msg string) error {
	rb.Lock()
	defer rb.Unlock()
	return rb.client.Publish(channel, msg).Err()
}

// Subscribe subscribes to a given redis channel.
// The callback function is triggered when a new notification is published
// to the subscribed channel.
func (rb *redisBroadcaster) Subscribe(channel string, f broadcaster.Callback) error {
	rb.Lock()
	defer rb.Unlock()
	if err := rb.pubSub.Subscribe(channel); err != nil {
		return err
	}
	rb.callback[channel] = f
	return nil
}

// UnSubscribe unsubscribe the client from a channel.
func (rb *redisBroadcaster) UnSubscribe(channel string) error {
	rb.Lock()
	defer rb.Unlock()
	if err := rb.pubSub.Unsubscribe(channel); err != nil {
		return err
	}
	delete(rb.callback, channel)
	return nil
}

// Close closes all PubSub streams and does the cleanup.
func (rb *redisBroadcaster) Close() error {
	rb.Lock()
	defer rb.Unlock()
	if err := rb.pubSub.Close(); err != nil {
		return err
	}
	rb.callback = make(map[string]broadcaster.Callback, rb.bufLen)
	rb.closeChan <- struct{}{}
	return nil
}
