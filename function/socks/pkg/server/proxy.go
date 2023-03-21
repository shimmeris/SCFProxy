package server

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/armon/go-socks5"
)

type Event struct {
	Key   string
	Addr  string
	Auth  string
	Stype string
}

type ScfClient interface {
	GetStream() (net.Conn, error)
}

func Handle(event Event) error {
	user, pass := "", ""

	userpass := strings.SplitN(event.Auth, ":", 2)
	if len(userpass) == 2 {
		user, pass = userpass[0], userpass[1]
	}

	socksServer := createSocks5(user, pass)
	for {
		scfClient, err := getScfClient(event.Addr, event.Key, event.Stype)
		if err != nil {
			continue
		}
		for {
			stream, err := scfClient.GetStream()
			if err != nil {
				break
			}
			go func() {
				err := socksServer.ServeConn(stream)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}
	}
}

func getScfClient(addr, key, scfType string) (ScfClient, error) {
	switch scfType {
	case "yamux":
		return NewYamuxScfClient(addr, key)
	case "quic":
		return NewQuicScfClient(addr, key)
	default:
		return nil, fmt.Errorf("Not this scf client type %s", scfType)
	}
}

func createSocks5(username, password string) *socks5.Server {
	conf := &socks5.Config{}
	if username == "" && password == "" {
		conf.AuthMethods = []socks5.Authenticator{socks5.NoAuthAuthenticator{}}
	} else {
		cred := socks5.StaticCredentials{username: password}
		conf.Credentials = cred
	}
	server, err := socks5.New(conf)
	if err != nil {
		log.Fatal(err)
	}
	return server
}
