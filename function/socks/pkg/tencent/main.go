package main

import (
	"github.com/tencentyun/scf-go-lib/cloudfunction"

	"socks/server"
)

func main() {
	cloudfunction.Start(server.Handle)
}
