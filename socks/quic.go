package socks

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	quic "github.com/quic-go/quic-go"
	"github.com/sirupsen/logrus"
)

type QuicScfServer struct {
	Conn quic.Connection
}

type QuicStreamWrapper struct {
	quic.Stream
}

func (s *QuicStreamWrapper) Close() error {
	s.Stream.CancelRead(0)
	return s.Stream.Close()
}

func (q *QuicScfServer) Listen(port, key string) {
	ln, err := quic.ListenAddr("0.0.0.0:"+port, generateTLSConfig(), nil)
	if err != nil {
		logrus.Fatal(err)
	}

	for {
		conn, err := ln.Accept(context.Background())
		if err != nil {
			logrus.Fatal(err)
		}
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			logrus.Fatal(err)
		}
		buf := make([]byte, KeyLength)
		stream.Read(buf)
		stream.CancelRead(0)
		stream.Close()
		if string(buf) == key {
			fmt.Printf("New scf connection from %s\n", conn.RemoteAddr().String())
			q.Conn = conn
		}
	}
}

func (q *QuicScfServer) GetStream() Stream {
	for {
		if q.Conn == nil {
			time.Sleep(5 * time.Second)
			continue
		}
		stream, err := q.Conn.OpenStream()
		if err != nil {
			logrus.Debug(err)
			q.Conn = nil
			continue
		}
		return &QuicStreamWrapper{Stream: stream}
	}
}

func (q *QuicScfServer) Close() error {
	return nil
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"socks"},
	}
}
