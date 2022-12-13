package main

import (
	"encoding/base64"
	"encoding/json"

	"github.com/aliyun/fc-runtime-go-sdk/events"
	"github.com/aliyun/fc-runtime-go-sdk/fc"

	"socks/server"
)

func HandleRequest(event events.TimerEvent) error {
	opts := &server.Options{}
	message, err := base64.StdEncoding.DecodeString(*event.Payload)
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
	fc.Start(HandleRequest)
}
