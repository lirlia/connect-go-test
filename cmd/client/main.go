package main

import (
	"context"
	"log"
	"net/http"

	hellov1 "example/gen/hello/v1"
	"example/gen/hello/v1/hellov1connect"
	"example/pkg/hash"

	"github.com/bufbuild/connect-go"
)

func main() {
	client := hellov1connect.NewHelloServiceClient(
		http.DefaultClient,
		"http://localhost:8081",
	)
	req := connect.NewRequest(&hellov1.HelloRequest{Name: "Jane"})
	h, err := hash.GenerateHash(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header().Set("request-hash", h)

	res, err := client.Hello(
		context.Background(),
		req,
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Hello)
}
