package server

import "testing"

func TestYamuxHandle(t *testing.T) {
	Handle(Event{
		Key:   "ung3oozaeTheil2a",
		Addr:  "127.0.0.1:10899",
		Auth:  "admin:admin",
		Stype: "yamux",
	})
}

func TestQuicHandle(t *testing.T) {
	Handle(Event{
		Key:   "ung3oozaeTheil2a",
		Addr:  "127.0.0.1:10899",
		Auth:  "admin:admin",
		Stype: "quic",
	})
}
