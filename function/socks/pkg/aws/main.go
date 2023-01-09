package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"socks/server"
)

func main() {
	lambda.Start(server.Handle)
}
