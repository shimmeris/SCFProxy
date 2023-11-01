package socks

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type YamuxScfServer struct {
	Sessions []*yamux.Session
}

func (y *YamuxScfServer) Listen(port, key string) {
	ln, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		logrus.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Fatal(err)
		}

		buf := make([]byte, KeyLength)
		conn.Read(buf)
		if string(buf) == key {
			fmt.Printf("New scf connection from %s\n", conn.RemoteAddr().String())
			session, err := yamux.Client(conn, nil)
			if err != nil {
				logrus.Error(err)
			}
			y.Sessions = append(y.Sessions, session)
		}
	}
}

func (y *YamuxScfServer) GetStream() Stream {
	for {
		l := len(y.Sessions)
		if l == 0 {
			logrus.Debug("No scf server connections")
			time.Sleep(5 * time.Second)
			continue
		}
		n := rand.Intn(l)
		conn, err := y.Sessions[n].Open()
		// remove inactive connections
		if err != nil {
			fmt.Printf("Remove invalid connection from %s\n", y.Sessions[n].RemoteAddr().String())
			y.Sessions[n].Close()
			y.Sessions = slices.Delete(y.Sessions, n, n+1)
			continue
		}
		return conn
	}
}

func (y *YamuxScfServer) Close() error {
	for _, s := range y.Sessions {
		s.Close()
	}
	return nil
}
