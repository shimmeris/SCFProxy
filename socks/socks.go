package socks

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

const KeyLength = 8

var sessions []*yamux.Session

func listenScf(port, key string) {
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
			fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())
			session, err := yamux.Client(conn, nil)
			if err != nil {
				logrus.Error(err)
			}
			sessions = append(sessions, session)
		}
	}
}

func listenClient(port string) {
	ln, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		logrus.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Error("listen client failed")
		}
		go forward(conn)
	}
}

func forward(conn net.Conn) {
	scfConn := pickConn(sessions)

	_forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}

	go _forward(conn, scfConn)
	go _forward(scfConn, conn)

}

func pickConn(sessions []*yamux.Session) net.Conn {
	for {
		l := len(sessions)
		if l == 0 {
			logrus.Debug("No scf server connections")
			time.Sleep(5 * time.Second)
			continue
		}
		n := rand.Intn(l)
		conn, err := sessions[n].Open()

		// remove inactive connections
		if err != nil {
			fmt.Printf("Remove invalid connection from %s\n", sessions[n].RemoteAddr().String())
			sessions[n].Close()
			sessions = slices.Delete(sessions, n, n+1)
			continue
		}

		return conn
	}
}

func Serve(socksPort, scfPort, key string) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			for _, s := range sessions {
				s.Close()
			}
			os.Exit(0)
		}()
	}()

	go listenScf(scfPort, key)
	listenClient(socksPort)

}
