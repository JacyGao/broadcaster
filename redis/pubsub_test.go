package redis

import (
	"github.com/go-redis/redis"
	"sync"
	"testing"
	"fmt"
	"time"
)

func TestPubSub(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	b, err := NewRedisBroadcaster(client,5)
	if err != nil {
		t.Fatal(err.Error())
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cs := fmt.Sprintf("%s%d", "test", i)
			if err := b.Subscribe(cs, func(s string) {
				fmt.Println(s)
			}); err != nil {
				t.Fatal(err.Error())
			}
		}(i)
	}
	wg.Wait()
	//fmt.Printf("%+v", b.pubSub)
	for i := 0; i < 5; i++ {
		cs := fmt.Sprintf("%s%d", "test", i)
		if err := b.Publish(cs, "ping"); err != nil {
			t.Fatal(err.Error())
		}
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cs := fmt.Sprintf("%s%d", "test", i)
			if err := b.UnSubscribe(cs); err != nil {
				t.Fatal(err.Error())
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < 5; i++ {
		cs := fmt.Sprintf("%s%d", "test", i)
		if err := b.Publish(cs, "ping"); err != nil {
			t.Fatal(err.Error())
		}
	}

	time.Sleep(100 * time.Millisecond)

	if err := b.pubSub.Close(); err != nil {
		t.Fatal(err.Error())
	}
	if err := client.Close(); err != nil {
		t.Fatal(err.Error())
	}
}