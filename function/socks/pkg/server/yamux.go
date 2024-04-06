package server

import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/yamux"
)

type YamuxScfClient struct {
	Session *yamux.Session
}

func NewYamuxScfClient(addr, key string) (*YamuxScfClient, error) {
	for i := 0; i < 5; i++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			if i == 4 {
				fmt.Printf("Connect to %s failed", conn.RemoteAddr().String())
			}
			time.Sleep(time.Duration((i+1)*5) * time.Second)
			fmt.Printf("[%d] Reconnecting\n", i)
			continue
		}
		conn.Write([]byte(key))
		session, err := yamux.Server(conn, nil)
		if err != nil {
			return nil, err
		}
		return &YamuxScfClient{Session: session}, nil
	}
	return nil, fmt.Errorf("Failed")
}

func (y *YamuxScfClient) GetStream() (net.Conn, error) {
	return y.Session.AcceptStream()
}
