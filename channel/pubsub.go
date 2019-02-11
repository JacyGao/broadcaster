package channel

import (
	"github.com/jacygao/broadcaster"
	"sync"
)

type message struct {
	channel string
	msg     string
}

type channelBroadcaster struct {
	sync.Mutex
	input     chan message
	outputs   map[string]broadcaster.Callback
	closeChan chan struct{}
}

func (b *channelBroadcaster) broadcast(msg message) {
	f, ok := b.outputs[msg.channel]
	if !ok {
		return
	}
	f(msg.msg)
}

func (b *channelBroadcaster) run() {
	for {
		select {
		case m := <-b.input:
			b.broadcast(m)

		case <-b.closeChan:
			return
		}
	}
}

// NewRedisBroadcaster returns an implementation of the Broadcaster interface backed by Go Channels.
func NewChannelBroadcaster() *channelBroadcaster {
	b := &channelBroadcaster{
		input:     make(chan message),
		outputs:   make(map[string]broadcaster.Callback),
		closeChan: make(chan struct{}, 1),
	}

	go b.run()

	return b
}

// Subscribe subscribes to a given channel.
// The callback function is triggered when a new notification is published
// to the subscribed channel.
func (b *channelBroadcaster) Subscribe(channel string, f broadcaster.Callback) error {
	b.Lock()
	defer b.Unlock()
	b.outputs[channel] = f
	return nil
}

// UnSubscribe unsubscribe the client from a channel.
func (b *channelBroadcaster) UnSubscribe(channel string) error {
	b.Lock()
	defer b.Unlock()
	delete(b.outputs, channel)
	return nil
}

// Close closes all PubSub streams and does the cleanup.
func (b *channelBroadcaster) Close() error {
	b.Lock()
	defer b.Unlock()
	close(b.input)
	b.outputs = make(map[string]broadcaster.Callback)
	b.closeChan <- struct{}{}
	return nil
}

// Publish publishes a new notification to a given channel.
func (b *channelBroadcaster) Publish(channel, msg string) error {
	if b != nil {
		b.input <- message{
			channel: channel,
			msg:     msg,
		}
	}
	return nil
}
