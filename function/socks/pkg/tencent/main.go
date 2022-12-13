package main

import (
	"encoding/base64"
	"encoding/json"

	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"

	"socks/server"
)

func handler(event events.TimerEvent) error {
	opts := &server.Options{}
	message, err := base64.StdEncoding.DecodeString(event.Message)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(message, opts); err != nil {
		return err
	}

	server.Run(opts)
	return nil
}

func main() {
	cloudfunction.Start(handler)
}
