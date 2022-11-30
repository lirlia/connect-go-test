package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	hellov1 "example/gen/hello/v1"        // generated by protoc-gen-go
	"example/gen/hello/v1/hellov1connect" // generated by protoc-gen-connect-go
	"example/pkg/hash"
)

type HelloServer struct{}

func (s *HelloServer) Hello(
	ctx context.Context,
	req *connect.Request[hellov1.HelloRequest],
) (*connect.Response[hellov1.HelloResponse], error) {

	res := connect.NewResponse(&hellov1.HelloResponse{
		Hello: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})

	return res, nil
}

func NewRequestHashCheckInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {

			expectedHash, err := hash.GenerateHash(req)
			if err != nil {
				return nil, err
			}
			// リクエストヘッダーからハッシュ値を取得
			requestHash := req.Header().Get("request-hash")

			// リクエストヘッダーのハッシュ値とリクエストのハッシュ値が一致しない場合はエラーを返す
			if requestHash != expectedHash {
				log.Printf("request hash mismatch expected: %s, actual: %s \n", requestHash, expectedHash)
				return nil, fmt.Errorf("request hash mismatch")
			}
			log.Printf("request hash matched! expected: %s, actual: %s \n", requestHash, expectedHash)
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func main() {

	mux := http.NewServeMux()
	interceptors := connect.WithInterceptors(NewRequestHashCheckInterceptor())
	path, handler := hellov1connect.NewHelloServiceHandler(&HelloServer{}, interceptors)

	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8081",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
