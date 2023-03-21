package socks

import (
	"testing"
)

func TestYamuxServe(t *testing.T) {
	Serve("10898", "10899", "ung3oozaeTheil2a", "yamux")
}

func TestQuicServe(t *testing.T) {
	Serve("10898", "10899", "ung3oozaeTheil2a", "quic")
}
