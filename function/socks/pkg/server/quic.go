package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	quic "github.com/quic-go/quic-go"
)

type QuicStreamWrapper struct {
	quic.Stream
	RAddr net.Addr
	LAddr net.Addr
}

func (q *QuicStreamWrapper) LocalAddr() net.Addr {
	return q.LAddr
}

func (q *QuicStreamWrapper) RemoteAddr() net.Addr {
	return q.RAddr
}

type QuicScfClient struct {
	Conn quic.Connection
}

func NewQuicScfClient(addr, key string) (*QuicScfClient, error) {
	for i := 0; i < 5; i++ {
		conn, err := quic.DialAddr(addr, &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"socks"},
		}, nil)
		if err != nil {
			if i == 4 {
				fmt.Printf("Connect to %s failed", conn.RemoteAddr().String())
			}
			time.Sleep(time.Duration((i+1)*5) * time.Second)
			fmt.Printf("[%d] Reconnecting\n", i)
			continue
		}
		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			time.Sleep(time.Duration((i+1)*5) * time.Second)
			fmt.Printf("[%d] Reconnecting\n", i)
		}
		stream.Write([]byte(key))
		stream.CancelRead(0)
		stream.Close()
		return &QuicScfClient{Conn: conn}, nil
	}
	return nil, fmt.Errorf("Not yet implemented")
}

func (q *QuicScfClient) GetStream() (net.Conn, error) {
	stream, err := q.Conn.AcceptStream(context.Background())
	return &QuicStreamWrapper{Stream: stream, RAddr: q.Conn.RemoteAddr(), LAddr: q.Conn.LocalAddr()}, err
}
