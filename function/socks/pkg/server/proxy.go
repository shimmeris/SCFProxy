package server

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/armon/go-socks5"
	"github.com/hashicorp/yamux"
)

type Options struct {
	Key  string
	Addr string
	Auth string
}

func Run(opts *Options) {
	user, pass := "", ""

	userpass := strings.SplitN(opts.Auth, ":", 2)
	if len(userpass) == 2 {
		user, pass = userpass[0], userpass[1]
	}

	socksServer := createSocks5(user, pass)
	for {
		conn := keepConnect(opts.Addr, opts.Key)
		session, err := yamux.Server(conn, nil)
		if err != nil {
			continue
		}
		for {
			stream, err := session.Accept()
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

func keepConnect(addr, key string) net.Conn {
	for i := 0; i < 5; i++ {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			time.Sleep(5 * time.Second)
			fmt.Println("Reconnecting")
			continue
		}
		conn.Write([]byte(key))
		return conn
	}
	return nil
}
