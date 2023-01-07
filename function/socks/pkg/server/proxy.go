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

type Event struct {
	Key  string
	Addr string
	Auth string
}

func Handle(event Event) error {
	user, pass := "", ""

	userpass := strings.SplitN(event.Auth, ":", 2)
	if len(userpass) == 2 {
		user, pass = userpass[0], userpass[1]
	}

	socksServer := createSocks5(user, pass)
	for {
		conn := keepConnect(event.Addr, event.Key)
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
			if i == 4 {
				fmt.Printf("Connect to %s failed", conn.RemoteAddr().String())
				return nil
			}
			time.Sleep(time.Duration((i+1)*5) * time.Second)
			fmt.Printf("[%d] Reconnecting\n", i)
			continue
		}
		conn.Write([]byte(key))
		return conn
	}
	return nil
}
