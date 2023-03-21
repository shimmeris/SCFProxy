package socks

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
)

const KeyLength = 16

type ScfServer interface {
	Listen(port, key string)
	GetStream() Stream
	io.Closer
}

type Stream interface {
	io.Reader
	io.Writer
	io.Closer
}

func listenClient(scfServer ScfServer, port string) {
	ln, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		logrus.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Error("listen client failed")
		}
		fmt.Printf("New socks connection from %s\n", conn.RemoteAddr().String())
		go Forward(scfServer, conn)
	}
}

func Forward(scfServer ScfServer, conn Stream) {
	scfConn := scfServer.GetStream()
	defer scfConn.Close()
	defer conn.Close()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go forward(wg, conn, scfConn)
	wg.Add(1)
	go forward(wg, scfConn, conn)
	wg.Wait()
}

func forward(wg *sync.WaitGroup, src, dest Stream) {
	defer wg.Done()
	io.Copy(src, dest)
}

func getScfServer(scfType string) (ScfServer, error) {
	switch scfType {
	case "yamux":
		return &YamuxScfServer{
			Sessions: make([]*yamux.Session, 0),
		}, nil
	case "quic":
		return &QuicScfServer{}, nil
	default:
		return nil, fmt.Errorf("Not this Scf Server Type %s", scfType)
	}
}

func Serve(socksPort, scfPort, key, scfType string) {
	scfServer, err := getScfServer(scfType)
	if err != nil {
		return
	}
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			scfServer.Close()
			os.Exit(0)
		}()
	}()

	fmt.Printf("scf listening on 0.0.0.0 %s\n", scfPort)
	go scfServer.Listen(scfPort, key)
	fmt.Printf("socks listening on 0.0.0.0 %s\n", socksPort)
	listenClient(scfServer, socksPort)
}
