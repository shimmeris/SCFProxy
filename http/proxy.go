package http

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/mitm"
	"github.com/sirupsen/logrus"
)

type Options struct {
	ListenAddr string
	CertPath   string
	KeyPath    string
	Apis       []string
}

func ServeProxy(opts *Options) error {
	p := martian.NewProxy()
	defer p.Close()

	// Prevent scfproxy from recursively connecting to itself.
	_, lport, _ := net.SplitHostPort(opts.ListenAddr)
	p.SetDial(func(network, address string) (net.Conn, error) {
		host, port, _ := net.SplitHostPort(address)
		if port == lport && (host == "localhost" || host == "127.0.0.1" || host == "::1") {
			return nil, errors.New("Detecting recursive connection")
		}
		return net.Dial(network, address)
	})

	l, err := net.Listen("tcp", opts.ListenAddr)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := configureTls(p, opts.CertPath, opts.KeyPath); err != nil {
		logrus.Error(err)
	}

	modifier, err := NewScfModifier(opts.Apis, lport)
	if err != nil {
		return err
	}

	p.SetRequestModifier(modifier)
	p.SetResponseModifier(modifier)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			os.Exit(0)
		}()
	}()

	fmt.Printf("HTTP proxy listening at %s\n", opts.ListenAddr)
	return p.Serve(l)
}

func configureTls(p *martian.Proxy, certPath, keyPath string) error {
	x509c, pk, err := GetX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}

	mc, err := mitm.NewConfig(x509c, pk)
	if err != nil {
		return err
	}

	p.SetMITM(mc)
	return nil

}
