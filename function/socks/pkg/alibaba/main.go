package main

import (
	"github.com/aliyun/fc-runtime-go-sdk/fc"

	"socks/server"
)

func main() {
	fc.Start(server.Handle)
}
