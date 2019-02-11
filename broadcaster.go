package broadcaster

type Callback func(s string)

type Broadcaster interface {
	Publish(channel, msg string) error
	Subscribe(channel string, f Callback) error
	UnSubscribe(channel string) error
	Close() error
}