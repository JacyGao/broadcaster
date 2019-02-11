package channel

import (
	"testing"
)

func TestBroadcast(t *testing.T) {
	mockVal := "mock_value"
	b := NewChannelBroadcaster()

	if err := b.Subscribe("test1", func(s string) {
		if s != mockVal {
			t.Fatalf("expected %s but got %s", mockVal, s)
		}
	}); err != nil {
		t.Fatal(err.Error())
	}

	if err := b.Subscribe("test2", func(s string) {
		if s != mockVal {
			t.Fatalf("expected %s but got %s", mockVal, s)
		}
	}); err != nil {
		t.Fatal(err.Error())
	}

	if len(b.outputs) != 2 {
		t.Fatalf("expected %d but got %d", 2, len(b.outputs))
	}

	if err := b.Publish("test1", mockVal); err != nil {
		t.Fatal(err.Error())
	}

	if err := b.Close(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestSubUnSub(t *testing.T) {
	mockVal := "mock_value"
	b := NewChannelBroadcaster()

	if err := b.Subscribe("test1", func(s string) {
		if s != mockVal {
			t.Fatalf("expected %s but got %s", mockVal, s)
		}
	}); err != nil {
		t.Fatal(err.Error())
	}

	if err := b.Subscribe("test2", func(s string) {
		if s != mockVal {
			t.Fatalf("expected %s but got %s", mockVal, s)
		}
	}); err != nil {
		t.Fatal(err.Error())
	}

	if err := b.UnSubscribe("test1"); err != nil {
		t.Fatalf(err.Error())
	}

	if len(b.outputs) != 1 {
		t.Fatalf("expected %d but got %d", 1, len(b.outputs))
	}
}